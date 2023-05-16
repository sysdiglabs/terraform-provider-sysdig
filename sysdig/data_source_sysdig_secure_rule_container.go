package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleContainer() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleContainerRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleDataSourceSchema(map[string]*schema.Schema{
			"matching": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"containers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func dataSourceSysdigRuleContainerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	ruleName := d.Get("name").(string)
	ruleType := v2.RuleTypeContainer

	rules, err := client.GetRuleGroup(ctx, ruleName, ruleType)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(rules) == 0 {
		return diag.Errorf("unable to find rule")
	}

	if len(rules) > 1 {
		return diag.Errorf("more than one rule with that name was found")
	}

	rule := rules[0]

	ruleDataSourceToResourceData(rule, d)

	_ = d.Set("matching", rule.Details.Containers.MatchItems)
	_ = d.Set("containers", rule.Details.Containers.Items)

	return nil
}
