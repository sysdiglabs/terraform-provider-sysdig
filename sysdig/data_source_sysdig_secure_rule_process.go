package sysdig

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Non-functional; see secure_rule_fast_engine_removed.go.
func dataSourceSysdigSecureRuleProcess() *schema.Resource {
	const typeName = "sysdig_secure_rule_process"
	const ruleType = "PROCESS"

	return &schema.Resource{
		DeprecationMessage: fastEngineDeprecation("data source", typeName, ruleType),
		ReadContext:        fastEngineDataSourceRemoved(typeName, ruleType),

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
