package sysdig

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Non-functional; see secure_rule_fast_engine_removed.go.
func dataSourceSysdigSecureRuleFilesystem() *schema.Resource {
	const typeName = "sysdig_secure_rule_filesystem"
	const ruleType = "FILESYSTEM"

	return &schema.Resource{
		DeprecationMessage: fastEngineDeprecation("data source", typeName, ruleType),
		ReadContext:        fastEngineDataSourceRemoved(typeName, ruleType),

		Schema: createRuleDataSourceSchema(map[string]*schema.Schema{
			"read_only": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"paths": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"read_write": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"paths": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		}),
	}
}
