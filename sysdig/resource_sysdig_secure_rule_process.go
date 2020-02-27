package sysdig

import (
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureRuleProcess() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigRuleProcessCreate,
		Update: resourceSysdigRuleProcessUpdate,
		Read:   resourceSysdigRuleProcessRead,
		Delete: resourceSysdigRuleProcessDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"matching": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"processes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func resourceSysdigRuleProcessCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	rule := resourceSysdigRuleProcessFromResourceData(d)

	rule, err := client.CreateRule(rule)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(rule.ID))
	d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleProcessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	rule, err := client.GetRuleByID(id)

	if err != nil {
		d.SetId("")
	}

	if rule.Details.Processes == nil {
		return errors.New("no process data for a process rule")
	}

	updateResourceDataForRule(d, rule)
	d.Set("matching", rule.Details.Processes.MatchItems)
	d.Set("processes", rule.Details.Processes.Items)

	return nil
}

func resourceSysdigRuleProcessUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	rule := resourceSysdigRuleProcessFromResourceData(d)

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err := client.UpdateRule(rule)

	return err
}

func resourceSysdigRuleProcessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return client.DeleteRule(id)
}

func resourceSysdigRuleProcessFromResourceData(d *schema.ResourceData) secure.Rule {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = "PROCESS"

	rule.Details.Processes = &secure.Processes{}
	rule.Details.Processes.MatchItems = d.Get("matching").(bool)
	rule.Details.Processes.Items = []string{}
	if processes, ok := d.Get("processes").([]interface{}); ok {
		for _, rawProcess := range processes {
			if process, ok := rawProcess.(string); ok {
				rule.Details.Processes.Items = append(rule.Details.Processes.Items, process)
			}
		}
	}

	return rule

}
