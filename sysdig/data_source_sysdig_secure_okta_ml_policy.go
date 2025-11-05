package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureOktaMLPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureOktaMLPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createOktaMLPolicyDataSourceSchema(),
	}
}

func dataSourceSysdigSecureOktaMLPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return oktaMLPolicyDataSourceRead(ctx, d, meta, "custom Okta ML policy", isCustomCompositePolicy)
}

func createOktaMLPolicyDataSourceSchema() map[string]*schema.Schema {
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
			MaxItems: 1,
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

func oktaMLPolicyDataSourceRead(ctx context.Context, d *schema.ResourceData, meta any, resourceName string, validationFunc func(v2.PolicyRulesComposite) bool) diag.Diagnostics {
	policy, err := compositePolicyDataSourceRead(ctx, d, meta, resourceName, policyTypeOktaML, validationFunc)
	if err != nil {
		return diag.FromErr(err)
	}

	err = oktaMLPolicyToResourceData(policy, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
