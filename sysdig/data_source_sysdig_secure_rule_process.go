package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleProcess() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleProcessRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleDataSourceSchema(map[string]*schema.Schema{
			"matching": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"processes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func dataSourceSysdigRuleProcessRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return commonDataSourceSysdigRuleRead(ctx, d, meta, v2.RuleTypeProcess, processRuleDataSourceToResourceData)
}

func processRuleDataSourceToResourceData(rule v2.Rule, d *schema.ResourceData) diag.Diagnostics {
	if rule.Details.Processes == nil {
		return diag.Errorf("no process data for a process rule")
	}
	_ = d.Set("matching", rule.Details.Processes.MatchItems)
	_ = d.Set("processes", rule.Details.Processes.Items)

	return nil
}
