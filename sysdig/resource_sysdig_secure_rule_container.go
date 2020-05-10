package sysdig

import (
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureRuleContainer() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigRuleContainerCreate,
		Update: resourceSysdigRuleContainerUpdate,
		Read:   resourceSysdigRuleContainerRead,
		Delete: resourceSysdigRuleContainerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"matching": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"containers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func resourceSysdigRuleContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	rule := resourceSysdigRuleContainerFromResourceData(d)

	rule, err = client.CreateRule(rule)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(rule.ID))
	d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleContainerRead(d *schema.ResourceData, meta interface{}) error {
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

	if rule.Details.Containers == nil {
		return errors.New("no container data for a container rule")
	}

	updateResourceDataForRule(d, rule)
	d.Set("matching", rule.Details.Containers.MatchItems)
	d.Set("containers", rule.Details.Containers.Items)

	return nil
}

func resourceSysdigRuleContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	rule := resourceSysdigRuleContainerFromResourceData(d)

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateRule(rule)

	return err
}

func resourceSysdigRuleContainerDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceSysdigRuleContainerFromResourceData(d *schema.ResourceData) secure.Rule {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = "CONTAINER"

	rule.Details.Containers = &secure.Containers{}
	rule.Details.Containers.MatchItems = d.Get("matching").(bool)
	rule.Details.Containers.Items = []string{}
	if containers, ok := d.Get("containers").([]interface{}); ok {
		for _, rawContainer := range containers {
			if container, ok := rawContainer.(string); ok {
				rule.Details.Containers.Items = append(rule.Details.Containers.Items, container)
			}
		}
	}

	return rule

}
