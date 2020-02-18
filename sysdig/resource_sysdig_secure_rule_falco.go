package sysdig

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureRuleFalco() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigRuleFalcoCreate,
		Update: resourceSysdigRuleFalcoUpdate,
		Read:   resourceSysdigRuleFalcoRead,
		Delete: resourceSysdigRuleFalcoDelete,

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
		}),
	}
}

func resourceSysdigRuleFalcoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	rule := resourceSysdigRuleFalcoFromResourceData(d)

	rule, err := client.CreateRule(rule)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(rule.ID))
	d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleFalcoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	rule, err := client.GetRuleByID(id)

	if err != nil {
		d.SetId("")
	}

	if rule.Details.Condition == nil {
		return errors.New("no condition data for a falco rule")
	}

	updateResourceDataForRule(d, rule)
	d.Set("condition", rule.Details.Condition.Condition)
	d.Set("output", rule.Details.Output)
	d.Set("priority", strings.ToLower(rule.Details.Priority))
	d.Set("source", rule.Details.Source)

	return nil
}

func resourceSysdigRuleFalcoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	rule := resourceSysdigRuleFalcoFromResourceData(d)

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err := client.UpdateRule(rule)

	return err
}

func resourceSysdigRuleFalcoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return client.DeleteRule(id)
}

func resourceSysdigRuleFalcoFromResourceData(d *schema.ResourceData) secure.Rule {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = "FALCO"

	rule.Details.Append = false
	rule.Details.Source = d.Get("source").(string)
	rule.Details.Output = d.Get("output").(string)
	rule.Details.Priority = d.Get("priority").(string)
	rule.Details.Condition = &secure.Condition{
		Condition:  d.Get("condition").(string),
		Components: []interface{}{},
	}

	return rule
}
