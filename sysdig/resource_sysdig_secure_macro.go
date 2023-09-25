package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"minimum_engine_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getSecureMacroClient(c SysdigClients) (v2.MacroInterface, error) {
	return c.sysdigSecureClientV2()
}

func resourceSysdigMacroCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureMacroClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	macro := macroFromResourceData(d)
	macro, err = client.CreateMacro(ctx, macro)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	d.SetId(strconv.Itoa(macro.ID))
	_ = d.Set("version", macro.Version)

	return nil
}

func resourceSysdigMacroUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureMacroClient(sysdigClients)
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
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigMacroRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureMacroClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	macro, err := client.GetMacroByID(ctx, id)

	if err != nil {
		d.SetId("")
	}

	_ = d.Set("name", macro.Name)
	_ = d.Set("version", macro.Version)
	_ = d.Set("condition", macro.Condition.Condition)
	_ = d.Set("append", macro.Append)
	if macro.MinimumEngineVersion != nil {
		_ = d.Set("minimum_engine_version", *macro.MinimumEngineVersion)
	}

	return nil
}

func resourceSysdigMacroDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureMacroClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteMacro(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func macroFromResourceData(d *schema.ResourceData) v2.Macro {
	macro := v2.Macro{
		Name:      d.Get("name").(string),
		Append:    d.Get("append").(bool),
		Condition: v2.MacroCondition{Condition: d.Get("condition").(string)},
	}
	minimumEngineVersionInterface, ok := d.GetOk("minimum_engine_version")
	if ok {
		minimumEngineVersion := minimumEngineVersionInterface.(int)
		macro.MinimumEngineVersion = &minimumEngineVersion
	}
	return macro
}
