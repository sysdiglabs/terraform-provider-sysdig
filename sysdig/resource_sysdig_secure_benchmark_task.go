package sysdig

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
)

func resourceSysdigSecureBenchmarkTask() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureBenchmarkTaskCreate,
		ReadContext:   resourceSysdigSecureBenchmarkTaskRead,
		UpdateContext: resourceSysdigSecureBenchmarkTaskUpdate,
		DeleteContext: resourceSysdigSecureBenchmarkTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schema": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(secure.SupportedBenchmarkTaskSchemas, false),
			},
			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSysdigSecureBenchmarkTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	benchmarkTask, err := client.CreateBenchmarkTask(ctx, benchmarkTaskFromResourceData(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(benchmarkTask.ID))
	_ = d.Set("name", benchmarkTask.Name)
	_ = d.Set("schema", benchmarkTask.Schema)
	_ = d.Set("scope", benchmarkTask.Scope)
	_ = d.Set("schedule", benchmarkTask.Schedule)
	_ = d.Set("enabled", benchmarkTask.Enabled)

	return nil
}

func resourceSysdigSecureBenchmarkTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	benchmarkTask, err := client.GetBenchmarkTask(ctx, d.Id())
	if err != nil {
		d.SetId("")
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(benchmarkTask.ID))
	_ = d.Set("name", benchmarkTask.Name)
	_ = d.Set("schema", benchmarkTask.Schema)
	_ = d.Set("scope", benchmarkTask.Scope)
	_ = d.Set("schedule", benchmarkTask.Schedule)
	_ = d.Set("enabled", benchmarkTask.Enabled)

	return nil
}

func resourceSysdigSecureBenchmarkTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	enabled := d.Get("enabled").(bool)

	if err := client.SetBenchmarkTaskEnabled(ctx, d.Id(), enabled); err != nil {
		d.SetId("")
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("enabled", enabled)

	return nil
}

func resourceSysdigSecureBenchmarkTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteBenchmarkTask(ctx, d.Id())
	if err != nil {
		d.SetId("")
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}

func benchmarkTaskFromResourceData(d *schema.ResourceData) *secure.BenchmarkTask {
	return &secure.BenchmarkTask{
		Name:     d.Get("name").(string),
		Schema:   d.Get("schema").(string),
		Scope:    d.Get("scope").(string),
		Schedule: d.Get("schedule").(string),
		Enabled:  d.Get("enabled").(bool),
	}
}
