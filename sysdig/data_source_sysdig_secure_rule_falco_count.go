package sysdig

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleFalcoCount() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleFalcoCountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateDiagFunc: validateDiagFunc(validateFalcoRuleSource),
			},
			"rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSysdigRuleFalcoCountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	ruleName := d.Get("name").(string)
	ruleType := d.Get("source").(string)
	rules, err := client.GetRuleGroup(ctx, ruleName, ruleType)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("count_%s", ruleName))
	_ = d.Set("name", ruleName)
	_ = d.Set("rule_count", len(rules))

	return nil
}
