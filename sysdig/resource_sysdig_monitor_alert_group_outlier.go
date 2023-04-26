package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigMonitorAlertGroupOutlier() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		DeprecationMessage: "Group Outlier Alerts have been deprecated, \"sysdig_monitor_alert_group_outlier\" will be removed in future releases",
		CreateContext:      resourceSysdigAlertGroupOutlierCreate,
		UpdateContext:      resourceSysdigAlertGroupOutlierUpdate,
		ReadContext:        resourceSysdigAlertGroupOutlierRead,
		DeleteContext:      resourceSysdigAlertGroupOutlierDelete,
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
			"monitor": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func resourceSysdigAlertGroupOutlierCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSysdigMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := groupOutlierAlertFromResourceData(data)
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

func resourceSysdigAlertGroupOutlierUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSysdigMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := groupOutlierAlertFromResourceData(data)
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

func resourceSysdigAlertGroupOutlierRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSysdigMonitorAlertClient(i.(SysdigClients))
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

	err = groupOutlierAlertToResourceData(&alert, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigAlertGroupOutlierDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSysdigMonitorAlertClient(i.(SysdigClients))
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

func groupOutlierAlertFromResourceData(data *schema.ResourceData) (alert *v2.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}

	alert.Type = "HOST_COMPARISON"

	for _, metric := range data.Get("monitor").([]interface{}) {
		alert.Monitor = append(alert.Monitor, &v2.Monitor{
			Metric:       metric.(string),
			StdDevFactor: 2,
		})
	}

	alert.SegmentCondition = &v2.SegmentCondition{Type: "ANY"}
	alert.SegmentBy = []string{"host.mac"}

	return
}

func groupOutlierAlertToResourceData(alert *v2.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	monitor_metrics := []string{}
	for _, v := range alert.Monitor {
		monitor_metrics = append(monitor_metrics, v.Metric)
	}
	_ = data.Set("monitor", monitor_metrics)

	return
}
