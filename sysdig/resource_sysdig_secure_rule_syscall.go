package sysdig

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Non-functional; see secure_rule_fast_engine_removed.go.
func resourceSysdigSecureRuleSyscall() *schema.Resource {
	const typeName = "sysdig_secure_rule_syscall"
	const ruleType = "SYSCALL"

	return &schema.Resource{
		DeprecationMessage: fastEngineDeprecation("resource", typeName, ruleType),

		CreateContext: fastEngineResourceRemoved(typeName, ruleType),
		ReadContext:   fastEngineRuleGone(typeName, ruleType),
		UpdateContext: fastEngineResourceRemoved(typeName, ruleType),
		DeleteContext: fastEngineRuleForget(typeName, ruleType),

		Schema: createRuleSchema(map[string]*schema.Schema{
			"matching": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"syscalls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}
