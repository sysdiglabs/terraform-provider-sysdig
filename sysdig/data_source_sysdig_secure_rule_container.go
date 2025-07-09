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

func dataSourceSysdigRuleContainerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return commonDataSourceSysdigRuleRead(ctx, d, meta, v2.RuleTypeContainer, containerRuleDataSourceToResourceData)
}

func containerRuleDataSourceToResourceData(rule v2.Rule, d *schema.ResourceData) diag.Diagnostics {
	_ = d.Set("matching", rule.Details.Containers.MatchItems)
	_ = d.Set("containers", rule.Details.Containers.Items)

	return nil
}
