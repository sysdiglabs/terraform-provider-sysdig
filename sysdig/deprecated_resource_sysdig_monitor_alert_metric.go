package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deprecatedResourceSysdigMonitorAlertMetric() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		DeprecationMessage: "\"sysdig_monitor_alert_metric\" has been deprecated and will be removed in future releases, use \"sysdig_monitor_alert_v2_metric\" instead",
		CreateContext:      deprecatedResourceSysdigAlertMetricCreate,
		UpdateContext:      deprecatedResourceSysdigAlertMetricUpdate,
		ReadContext:        deprecatedResourceSysdigAlertMetricRead,
		DeleteContext:      deprecatedResourceSysdigAlertMetricDelete,
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

func deprecatedResourceSysdigAlertMetricCreate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := deprecatedMetricAlertFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	alertCreated, err := client.CreateAlert(ctx, *alert)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(alertCreated.ID))
	_ = data.Set("version", alertCreated.Version)

	return nil
}

func deprecatedResourceSysdigAlertMetricUpdate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := deprecatedMetricAlertFromResourceData(data)
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

func deprecatedResourceSysdigAlertMetricRead(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := client.GetAlertByID(ctx, id)
	if err != nil {
		data.SetId("")
		return nil
	}

	err = deprecatedMetricAlertToResourceData(&alert, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deprecatedResourceSysdigAlertMetricDelete(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlertByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deprecatedMetricAlertFromResourceData(data *schema.ResourceData) (alert *v2.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}
	alert.Condition = data.Get("metric").(string)

	if alertsBy, ok := data.GetOk("multiple_alerts_by"); ok {
		alert.SegmentCondition = &v2.SegmentCondition{Type: "ANY"}
		for _, v := range alertsBy.([]any) {
			alert.SegmentBy = append(alert.SegmentBy, v.(string))
		}
	}
	return
}

func deprecatedMetricAlertToResourceData(alert *v2.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	_ = data.Set("metric", alert.Condition)
	_ = data.Set("multiple_alerts_by", alert.SegmentBy)

	return
}
