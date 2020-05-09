package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigSecureMacro() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigMacroCreate,
		Update: resourceSysdigMacroUpdate,
		Read:   resourceSysdigMacroRead,
		Delete: resourceSysdigMacroDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"append": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"condition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSysdigMacroCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	macro := macroFromResourceData(d)
	macro, err := client.CreateMacro(macro)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(macro.ID))
	d.Set("version", macro.Version)

	return nil
}

func resourceSysdigMacroUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	macro := macroFromResourceData(d)
	macro.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	macro.ID = id

	_, err := client.UpdateMacro(macro)
	return err
}

func resourceSysdigMacroRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	id, _ := strconv.Atoi(d.Id())
	macro, err := client.GetMacroById(id)

	if err != nil {
		d.SetId("")
	}

	d.Set("name", macro.Name)
	d.Set("version", macro.Version)
	d.Set("items", macro.Condition.Condition)
	d.Set("append", macro.Append)

	return nil
}

func resourceSysdigMacroDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteMacro(id)
}

func macroFromResourceData(d *schema.ResourceData) secure.Macro {
	return secure.Macro{
		Name:      d.Get("name").(string),
		Append:    d.Get("append").(bool),
		Condition: secure.MacroCondition{Condition: d.Get("condition").(string)},
	}
}
