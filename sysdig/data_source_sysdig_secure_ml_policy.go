package sysdig

import (
	"context"
	"fmt"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureMLPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureMLPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createMLPolicyDataSourceSchema(),
	}
}

func dataSourceSysdigSecureMLPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return mlPolicyDataSourceRead(ctx, d, meta, "custom ML policy", isCustomCompositePolicy)
}

func createMLPolicyDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// IMPORTANT: Type is implicit: It's automatically added upon conversion to JSON
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name":                  NameSchema(),
		"description":           DescriptionComputedSchema(),
		"enabled":               EnabledComputedSchema(),
		"severity":              SeverityComputedSchema(),
		"scope":                 ScopeComputedSchema(),
		"version":               VersionSchema(),
		"notification_channels": NotificationChannelsComputedSchema(),
		"runbook":               RunbookComputedSchema(),
		"rule": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id":                   ReadOnlyIntSchema(),
					"name":                 ReadOnlyStringSchema(),
					"description":          DescriptionComputedSchema(),
					"tags":                 TagsSchema(),
					"version":              VersionSchema(),
					"cryptomining_trigger": MLRuleThresholdAndSeverityComputedSchema(),
				},
			},
		},
	}
}

func mlPolicyDataSourceRead(ctx context.Context, d *schema.ResourceData, meta any, resourceName string, validationFunc func(v2.PolicyRulesComposite) bool) diag.Diagnostics {
	policy, err := compositePolicyDataSourceRead(ctx, d, meta, resourceName, policyTypeML, validationFunc)
	if err != nil {
		return diag.FromErr(err)
	}
	err = mlPolicyToResourceData(policy, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func compositePolicyDataSourceRead(ctx context.Context, d *schema.ResourceData, meta any, resourceName string, policyType string, validationFunc func(v2.PolicyRulesComposite) bool) (*v2.PolicyRulesComposite, error) {
	client, err := getSecureCompositePolicyClient(meta.(SysdigClients))
	if err != nil {
		return nil, err
	}

	policyName := d.Get("name").(string)

	policies, _, err := client.ListCompositePoliciesByNameAndType(ctx, policyType, policyName)
	if err != nil {
		return nil, err
	}

	var policy v2.PolicyRulesComposite
	for _, existingPolicy := range policies {
		tflog.Debug(ctx, "Filtered policies", map[string]any{"name": existingPolicy.Policy.Name})

		if existingPolicy.Policy.Name == policyName && existingPolicy.Policy.Type == policyType {
			if !validationFunc(existingPolicy) {
				return nil, fmt.Errorf("policy is not a %s", resourceName)
			}
			policy = existingPolicy
			break
		}
	}

	if policy.Policy == nil {
		return nil, fmt.Errorf("unable to find policy %s", resourceName)
	}

	if policy.Policy.ID == 0 {
		return nil, fmt.Errorf("unable to find %s", resourceName)
	}

	return &policy, nil
}
