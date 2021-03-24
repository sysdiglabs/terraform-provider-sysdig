package sysdig

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
)

func resourceSysdigSecureRuleFalco() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigRuleFalcoCreate,
		UpdateContext: resourceSysdigRuleFalcoUpdate,
		ReadContext:   resourceSysdigRuleFalcoRead,
		DeleteContext: resourceSysdigRuleFalcoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"condition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"priority": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "warning",
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"}, false)),
			},
			"source": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"syscall", "k8s_audit", "aws_cloudtrail"}, false)),
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

	rule, err := resourceSysdigRuleFalcoFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

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

	rule, err := resourceSysdigRuleFalcoFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

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

func resourceSysdigRuleFalcoFromResourceData(d *schema.ResourceData) (secure.Rule, error) {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = "FALCO"

	appendMode, appendModeIsSet := d.GetOk("append")
	if appendModeIsSet {
		ptr := appendMode.(bool)
		rule.Details.Append = &ptr
	}

	if source, ok := d.GetOk("source"); ok && source.(string) != "" {
		rule.Details.Source = source.(string)
	} else if !appendModeIsSet || !(appendMode.(bool)) {
		return secure.Rule{}, errors.New("source must be set when append = false")
	}

	if output, ok := d.GetOk("output"); ok && output.(string) != "" {
		rule.Details.Output = output.(string)
	} else if !appendModeIsSet || !(appendMode.(bool)) {
		return secure.Rule{}, errors.New("output must be set when append = false")
	}

	if priority, ok := d.GetOk("priority"); ok && priority.(string) != "" {
		rule.Details.Priority = priority.(string)
	} else if !appendModeIsSet || !(appendMode.(bool)) {
		return secure.Rule{}, errors.New("priority must be set when append = false")
	}

	rule.Details.Condition = &secure.Condition{
		Condition:  d.Get("condition").(string),
		Components: []interface{}{},
	}

	return rule, nil
}
