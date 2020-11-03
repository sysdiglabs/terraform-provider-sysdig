package sysdig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureRuleFalco() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigRuleFalcoCreate,
		UpdateContext: resourceSysdigRuleFalcoUpdate,
		ReadContext:   resourceSysdigRuleFalcoRead,
		DeleteContext: resourceSysdigRuleFalcoDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"condition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"emergency", "alert", "critical", "error", "warning", "notice", "informational", "informational", "debug"}, false),
			},
			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"syscall", "k8s_audit"}, false),
			},
			"append": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		}),
	}
}

func resourceSysdigRuleFalcoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	rule := resourceSysdigRuleFalcoFromResourceData(d)

	rule, err = client.CreateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(rule.ID))
	d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleFalcoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := client.GetRuleByID(ctx, id)

	if err != nil {
		d.SetId("")
	}

	if rule.Details.Condition == nil {
		return diag.Errorf("no condition data for a falco rule")
	}

	updateResourceDataForRule(d, rule)
	d.Set("condition", rule.Details.Condition.Condition)
	d.Set("output", rule.Details.Output)
	d.Set("priority", strings.ToLower(rule.Details.Priority))
	d.Set("source", rule.Details.Source)
	if rule.Details.Append != nil {
		d.Set("append", *rule.Details.Append)
	}

	return nil
}

func resourceSysdigRuleFalcoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	rule := resourceSysdigRuleFalcoFromResourceData(d)

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigRuleFalcoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteRule(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSysdigRuleFalcoFromResourceData(d *schema.ResourceData) secure.Rule {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = "FALCO"

	rule.Details.Source = d.Get("source").(string)
	rule.Details.Output = d.Get("output").(string)
	rule.Details.Priority = d.Get("priority").(string)
	rule.Details.Condition = &secure.Condition{
		Condition:  d.Get("condition").(string),
		Components: []interface{}{},
	}

	if appendMode, ok := d.GetOk("append"); ok {
		ptr := appendMode.(bool)
		rule.Details.Append = &ptr
	}

	return rule
}
