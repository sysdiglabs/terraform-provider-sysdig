package sysdig

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Non-functional; see secure_rule_fast_engine_removed.go.
func dataSourceSysdigSecureRuleNetwork() *schema.Resource {
	const typeName = "sysdig_secure_rule_network"
	const ruleType = "NETWORK"

	return &schema.Resource{
		DeprecationMessage: fastEngineDeprecation("data source", typeName, ruleType),
		ReadContext:        fastEngineDataSourceRemoved(typeName, ruleType),

		Schema: createRuleDataSourceSchema(map[string]*schema.Schema{
			"block_inbound": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"block_outbound": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tcp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"udp": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		}),
	}
}
