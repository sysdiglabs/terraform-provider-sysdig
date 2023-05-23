package sysdig

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleFalco() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleFalcoRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleDataSourceSchema(map[string]*schema.Schema{
			"source": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateDiagFunc: validateDiagFunc(validateFalcoRuleSource),
			},
			"index": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"condition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"append": {

				Type:     schema.TypeBool,
				Computed: true,
			},
			"exceptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comps": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"values": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		}),
	}
}

func dataSourceSysdigRuleFalcoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	ruleName := d.Get("name").(string)
	ruleType := d.Get("source").(string)
	ruleIndex := d.Get("index").(int)
	rules, err := client.GetRuleGroup(ctx, ruleName, ruleType)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(rules) == 0 {
		return diag.Errorf("unable to find rule")
	}

	if ruleIndex >= len(rules) {
		return diag.Errorf("unable to find rule at the index provided")
	}
	rule := rules[ruleIndex]

	ruleDataSourceToResourceData(rule, d)

	if rule.Details.Condition != nil {
		_ = d.Set("condition", rule.Details.Condition.Condition)
	}
	_ = d.Set("output", rule.Details.Output)
	_ = d.Set("priority", rule.Details.Priority)
	_ = d.Set("source", rule.Details.Source)
	if rule.Details.Append != nil {
		_ = d.Set("append", *rule.Details.Append)
	}
	if err := updateResourceDataExceptions(d, rule.Details.Exceptions); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
