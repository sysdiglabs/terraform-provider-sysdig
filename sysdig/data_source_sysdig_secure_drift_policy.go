package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureDriftPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureDriftPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createDriftPolicyDataSourceSchema(),
	}
}

func dataSourceSysdigSecureDriftPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return driftPolicyDataSourceRead(ctx, d, meta, "custom policy", isCustomCompositePolicy)
}

func createDriftPolicyDataSourceSchema() map[string]*schema.Schema {
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
		"rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id":          ReadOnlyIntSchema(),
					"name":        ReadOnlyStringSchema(),
					"description": DescriptionComputedSchema(),
					"tags":        TagsSchema(),
					"version":     VersionSchema(),
					"details": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"enabled":             BoolComputedSchema(),
								"exceptions":          ExceptionsComputedSchema(),
								"prohibited_binaries": ExceptionsComputedSchema(),
							},
						},
					},
				},
			},
		},
		"actions": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"prevent_drift": PreventActionComputedSchema(),
					"container":     ContainerActionComputedSchema(),
					"capture":       CaptureActionComputedSchema(),
				},
			},
		},
	}
}

func driftPolicyDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}, resourceName string, validationFunc func(v2.PolicyRulesComposite) bool) diag.Diagnostics {
	client, err := getSecureCompositePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policyName := d.Get("name").(string)
	policyType := "drift" // assume it's "drift" type

	policies, _, err := client.FilterCompositePoliciesByNameAndType(ctx, policyType, policyName)
	if err != nil {
		return diag.FromErr(err)
	}

	var policy v2.PolicyRulesComposite
	for _, existingPolicy := range policies {
		tflog.Debug(ctx, "Filtered policies", map[string]interface{}{"name": existingPolicy.Policy.Name})

		if existingPolicy.Policy.Name == policyName && existingPolicy.Policy.Type == policyType {
			if !validationFunc(existingPolicy) {
				return diag.Errorf("policy is not a %s", resourceName)
			}
			policy = existingPolicy
			break
		}
	}

	if policy.Policy == nil {
		return diag.Errorf("unable to find policy %s", resourceName)
	}

	if policy.Policy.ID == 0 {
		return diag.Errorf("unable to find %s", resourceName)
	}

	err = driftPolicyToResourceData(&policy, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
