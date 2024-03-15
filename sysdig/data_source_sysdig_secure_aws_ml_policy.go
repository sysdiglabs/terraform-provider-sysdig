package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureAWSMLPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureAWSMLPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createAWSMLPolicyDataSourceSchema(),
	}
}

func dataSourceSysdigSecureAWSMLPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return awsMLPolicyDataSourceRead(ctx, d, meta, "custom AWS ML policy", isCustomCompositePolicy)
}

func createAWSMLPolicyDataSourceSchema() map[string]*schema.Schema {
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
					"id":                      ReadOnlyIntSchema(),
					"name":                    ReadOnlyStringSchema(),
					"description":             DescriptionComputedSchema(),
					"tags":                    TagsSchema(),
					"version":                 VersionSchema(),
					"anomalous_console_login": MLRuleThresholdAndSeverityComputedSchema(),
				},
			},
		},
	}
}

func awsMLPolicyDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}, resourceName string, validationFunc func(v2.PolicyRulesComposite) bool) diag.Diagnostics {
	policy, err := compositePolicyDataSourceRead(ctx, d, meta, resourceName, policyTypeAWSML, validationFunc)
	if err != nil {
		return diag.FromErr(err)
	}

	err = awsMLPolicyToResourceData(policy, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
