package sysdig

import (
	"context"
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigSecureMacro() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMacroCreate,
		UpdateContext: resourceSysdigMacroUpdate,
		ReadContext:   resourceSysdigMacroRead,
		DeleteContext: resourceSysdigMacroDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

func resourceSysdigMacroCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	macro := macroFromResourceData(d)
	macro, err = client.CreateMacro(ctx, macro)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(macro.ID))
	d.Set("version", macro.Version)

	return nil
}

func resourceSysdigMacroUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	macro := macroFromResourceData(d)
	macro.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	macro.ID = id

	_, err = client.UpdateMacro(ctx, macro)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSysdigMacroRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	macro, err := client.GetMacroById(ctx, id)

	if err != nil {
		d.SetId("")
	}

	d.Set("name", macro.Name)
	d.Set("version", macro.Version)
	d.Set("condition", macro.Condition.Condition)
	d.Set("append", macro.Append)

	return nil
}

func resourceSysdigMacroDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteMacro(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func macroFromResourceData(d *schema.ResourceData) secure.Macro {
	return secure.Macro{
		Name:      d.Get("name").(string),
		Append:    d.Get("append").(bool),
		Condition: secure.MacroCondition{Condition: d.Get("condition").(string)},
	}
}
