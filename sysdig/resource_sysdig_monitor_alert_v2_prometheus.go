package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigMonitorAlertV2Prometheus() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorAlertV2PrometheusCreate,
		UpdateContext: resourceSysdigMonitorAlertV2PrometheusUpdate,
		ReadContext:   resourceSysdigMonitorAlertV2PrometheusRead,
		DeleteContext: resourceSysdigMonitorAlertV2PrometheusDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createAlertV2Schema(map[string]*schema.Schema{
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func getAlertV2PrometheusClient(c SysdigClients) (v2.AlertV2PrometheusInterface, error) {
	return c.sysdigMonitorClientV2()
}

func resourceSysdigMonitorAlertV2PrometheusCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2PrometheusClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a := buildAlertV2PrometheusStruct(d)

	aCreated, err := client.CreateAlertV2Prometheus(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(aCreated.ID))

	err = updateAlertV2PrometheusState(d, &aCreated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2PrometheusRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2PrometheusClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := client.GetAlertV2Prometheus(ctx, id)

	if err != nil {
		d.SetId("")
		return nil
	}

	err = updateAlertV2PrometheusState(d, &a)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2PrometheusUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2PrometheusClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a := buildAlertV2PrometheusStruct(d)

	a.ID, _ = strconv.Atoi(d.Id())

	aUpdated, err := client.UpdateAlertV2Prometheus(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateAlertV2PrometheusState(d, &aUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2PrometheusDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2PrometheusClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlertV2Prometheus(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildAlertV2PrometheusStruct(d *schema.ResourceData) *v2.AlertV2Prometheus {
	alertV2Common := alertV2CommonAdapterOldToNew(buildAlertV2CommonStruct(d))

	alertV2Common.Type = string(v2.Prometheus)

	config := v2.AlertV2ConfigPrometheus{}
	config.Query = d.Get("query").(string)

	alert := &v2.AlertV2Prometheus{
		AlertV2Common: *alertV2Common,
		Config:        config,
	}
	return alert
}

func updateAlertV2PrometheusState(d *schema.ResourceData, alert *v2.AlertV2Prometheus) (err error) {
	err = updateAlertV2CommonState(d, alertV2CommonAdapterNewToOld(&alert.AlertV2Common))
	if err != nil {
		return
	}

	_ = d.Set("query", alert.Config.Query)

	return
}
