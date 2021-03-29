package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
)

func resourceSysdigMonitorAlertMetric() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigAlertMetricCreate,
		UpdateContext: resourceSysdigAlertMetricUpdate,
		ReadContext:   resourceSysdigAlertMetricRead,
		DeleteContext: resourceSysdigAlertMetricDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createAlertSchema(map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"multiple_alerts_by": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func resourceSysdigAlertMetricCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := metricAlertFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	alertCreated, err := client.CreateAlert(ctx, *alert)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(alertCreated.ID))
	data.Set("version", alertCreated.Version)
	return nil
}

func resourceSysdigAlertMetricUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := metricAlertFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	alert.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateAlert(ctx, *alert)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigAlertMetricRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := client.GetAlertById(ctx, id)

	if err != nil {
		data.SetId("")
		return nil
	}

	err = metricAlertToResourceData(&alert, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigAlertMetricDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlert(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func metricAlertFromResourceData(data *schema.ResourceData) (alert *monitor.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}
	alert.Condition = data.Get("metric").(string)

	if alerts_by, ok := data.GetOk("multiple_alerts_by"); ok {
		alert.SegmentCondition = &monitor.SegmentCondition{Type: "ANY"}
		for _, v := range alerts_by.([]interface{}) {
			alert.SegmentBy = append(alert.SegmentBy, v.(string))
		}
	}
	return
}

func metricAlertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	data.Set("metric", alert.Condition)
	data.Set("multiple_alerts_by", alert.SegmentBy)

	return
}
