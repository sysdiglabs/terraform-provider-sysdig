package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
)

func deprecatedResourceSysdigMonitorAlertDowntime() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		DeprecationMessage: "\"sysdig_monitor_alert_downtime\" has been deprecated and will be removed in future releases, use \"sysdig_monitor_alert_v2_downtime\" instead",
		CreateContext:      deprecatedResourceSysdigAlertDowntimeCreate,
		UpdateContext:      deprecatedResourceSysdigAlertDowntimeUpdate,
		ReadContext:        deprecatedResourceSysdigAlertDowntimeRead,
		DeleteContext:      deprecatedResourceSysdigAlertDowntimeDelete,
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
			"entities_to_monitor": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"trigger_after_pct": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
		}),
	}
}

func deprecatedResourceSysdigAlertDowntimeCreate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := deprecatedDowntimeAlertFromResourceData(data)
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

func deprecatedResourceSysdigAlertDowntimeUpdate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := deprecatedDowntimeAlertFromResourceData(data)
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

func deprecatedResourceSysdigAlertDowntimeRead(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
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

	err = deprecatedDowntimeAlertToResourceData(&alert, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deprecatedResourceSysdigAlertDowntimeDelete(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
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

func deprecatedDowntimeAlertFromResourceData(d *schema.ResourceData) (alert *v2.Alert, err error) {
	alert, err = alertFromResourceData(d)
	if err != nil {
		return
	}

	alert.SegmentCondition = &v2.SegmentCondition{Type: "ANY"}
	alert.Condition = fmt.Sprintf("avg(timeAvg(uptime)) <= %.2f", 1.0-(cast.ToFloat64(d.Get("trigger_after_pct"))/100.0))

	entitiesRaw := d.Get("entities_to_monitor").([]any)
	for _, entityRaw := range entitiesRaw {
		alert.SegmentBy = append(alert.SegmentBy, entityRaw.(string))
	}

	return
}

func deprecatedDowntimeAlertToResourceData(alert *v2.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	var triggerAfterPct float64
	_, _ = fmt.Sscanf(alert.Condition, "avg(timeAvg(uptime)) <= %f", &triggerAfterPct)
	triggerAfterPct = (1 - triggerAfterPct) * 100

	_ = data.Set("trigger_after_pct", int(triggerAfterPct))
	_ = data.Set("entities_to_monitor", alert.SegmentBy)

	return
}
