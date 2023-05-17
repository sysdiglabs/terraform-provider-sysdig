package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleFileSystem() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleFileSystemRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleDataSourceSchema(map[string]*schema.Schema{
			"read_only": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Computed: true,
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
							Computed: true,
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

func dataSourceSysdigRuleFileSystemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	ruleName := d.Get("name").(string)
	ruleType := v2.RuleTypeFileSystem

	rules, err := client.GetRuleGroup(ctx, ruleName, ruleType)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(rules) == 0 {
		return diag.Errorf("unable to find rule")
	}

	if len(rules) > 1 {
		return diag.Errorf("more than one rule with that name was found")
	}

	rule := rules[0]

	ruleDataSourceToResourceData(rule, d)

	_ = d.Set("read_only", rule.Details.ReadPaths)
	_ = d.Set("read_write", rule.Details.ReadWritePaths)

	return nil
}
