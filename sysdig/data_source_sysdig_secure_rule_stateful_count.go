package sysdig

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleStatefulCount() *schema.Resource {
	timeout := 1 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleStatefulCountRead,

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
				Required:         true,
				ValidateDiagFunc: validateDiagFunc(validateStatefulRuleSource),
			},
			"rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSysdigRuleStatefulCountRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	ruleName := d.Get("name").(string)
	ruleType := d.Get("source").(string)
	rules, err := client.GetStatefulRuleGroup(ctx, ruleName, ruleType)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("count__%s__%s", ruleName, ruleType))
	_ = d.Set("name", ruleName)
	_ = d.Set("rule_count", len(rules))

	return nil
}
