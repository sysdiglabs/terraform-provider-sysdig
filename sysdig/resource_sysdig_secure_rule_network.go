package sysdig

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Non-functional; see secure_rule_fast_engine_removed.go.
func resourceSysdigSecureRuleNetwork() *schema.Resource {
	const typeName = "sysdig_secure_rule_network"
	const ruleType = "NETWORK"

	return &schema.Resource{
		DeprecationMessage: fastEngineDeprecation("resource", typeName, ruleType),

		CreateContext: fastEngineResourceRemoved(typeName, ruleType),
		ReadContext:   fastEngineRuleGone(typeName, ruleType),
		UpdateContext: fastEngineResourceRemoved(typeName, ruleType),
		DeleteContext: fastEngineRuleForget(typeName, ruleType),

		Schema: createRuleSchema(map[string]*schema.Schema{
			"block_inbound": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"block_outbound": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"tcp": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Required: true,
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
							Optional: true,
							Default:  true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Required: true,
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
