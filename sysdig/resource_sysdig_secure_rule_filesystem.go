package sysdig

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Non-functional; see secure_rule_fast_engine_removed.go.
func resourceSysdigSecureRuleFilesystem() *schema.Resource {
	const typeName = "sysdig_secure_rule_filesystem"
	const ruleType = "FILESYSTEM"

	return &schema.Resource{
		DeprecationMessage: fastEngineDeprecation("resource", typeName, ruleType),

		CreateContext: fastEngineResourceRemoved(typeName, ruleType),
		ReadContext:   fastEngineRuleGone(typeName, ruleType),
		UpdateContext: fastEngineResourceRemoved(typeName, ruleType),
		DeleteContext: fastEngineRuleForget(typeName, ruleType),

		Schema: createRuleSchema(map[string]*schema.Schema{
			"read_only": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"paths": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"read_write": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"paths": {
							Type:     schema.TypeList,
							Required: true,
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
