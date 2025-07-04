package sysdig

import (
	"context"
	"strconv"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Creates a schema with a default schema that a Secure Rule data source should have
// Additional fields will be passed in via the parameter
func createRuleDataSourceSchema(original map[string]*schema.Schema) map[string]*schema.Schema {
	ruleSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}

	for k, v := range original {
		ruleSchema[k] = v
	}

	return ruleSchema
}

func ruleDataSourceToResourceData(rule v2.Rule, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(rule.ID))

	_ = d.Set("name", rule.Name)
	_ = d.Set("description", rule.Description)
	_ = d.Set("tags", rule.Tags)
	_ = d.Set("version", rule.Version)
}

func commonDataSourceSysdigRuleRead(ctx context.Context, d *schema.ResourceData, meta any, ruleType string, setResourceData func(v2.Rule, *schema.ResourceData) diag.Diagnostics) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	ruleName := d.Get("name").(string)

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
	diag := setResourceData(rule, d)

	return diag
}
