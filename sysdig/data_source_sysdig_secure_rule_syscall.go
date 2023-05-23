package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleSyscall() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleSyscallRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleDataSourceSchema(map[string]*schema.Schema{
			"matching": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"syscalls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func dataSourceSysdigRuleSyscallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return commonDataSourceSysdigRuleRead(ctx, d, meta, v2.RuleTypeSyscall, syscallRuleDataSourceToResourceData)
}

func syscallRuleDataSourceToResourceData(rule v2.Rule, d *schema.ResourceData) diag.Diagnostics {
	if rule.Details.Syscalls == nil {
		return diag.Errorf("no syscall data for a syscall rule")
	}

	_ = d.Set("matching", rule.Details.Syscalls.MatchItems)
	_ = d.Set("syscalls", rule.Details.Syscalls.Items)

	return nil
}
