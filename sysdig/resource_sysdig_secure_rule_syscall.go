package sysdig

import (
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureRuleSyscall() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigRuleSyscallCreate,
		Update: resourceSysdigRuleSyscallUpdate,
		Read:   resourceSysdigRuleSyscallRead,
		Delete: resourceSysdigRuleSyscallDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

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

func resourceSysdigRuleSyscallCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	rule := resourceSysdigRuleSyscallFromResourceData(d)

	rule, err = client.CreateRule(rule)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(rule.ID))
	d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleSyscallRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	rule, err := client.GetRuleByID(id)

	if err != nil {
		d.SetId("")
	}

	if rule.Details.Syscalls == nil {
		return errors.New("no syscall data for a syscall rule")
	}

	updateResourceDataForRule(d, rule)
	d.Set("matching", rule.Details.Syscalls.MatchItems)
	d.Set("syscalls", rule.Details.Syscalls.Items)

	return nil
}

func resourceSysdigRuleSyscallUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	rule := resourceSysdigRuleSyscallFromResourceData(d)

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateRule(rule)

	return err
}

func resourceSysdigRuleSyscallDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return client.DeleteRule(id)
}

func resourceSysdigRuleSyscallFromResourceData(d *schema.ResourceData) secure.Rule {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = "SYSCALL"

	rule.Details.Syscalls = &secure.Syscalls{}
	rule.Details.Syscalls.MatchItems = d.Get("matching").(bool)
	rule.Details.Syscalls.Items = []string{}
	if syscalls, ok := d.Get("syscalls").([]interface{}); ok {
		for _, rawSyscall := range syscalls {
			if syscall, ok := rawSyscall.(string); ok {
				rule.Details.Syscalls.Items = append(rule.Details.Syscalls.Items, syscall)
			}
		}
	}

	return rule

}
