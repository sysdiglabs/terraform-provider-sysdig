package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func deprecatedResourceSysdigMonitorAlertGroupOutlier() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		DeprecationMessage: "Group Outlier Alerts have been deprecated, \"sysdig_monitor_alert_group_outlier\" will be removed in future releases",
		CreateContext:      deprecatedResourceSysdigAlertGroupOutlierCreate,
		UpdateContext:      deprecatedResourceSysdigAlertGroupOutlierUpdate,
		ReadContext:        deprecatedResourceSysdigAlertGroupOutlierRead,
		DeleteContext:      deprecatedResourceSysdigAlertGroupOutlierDelete,
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

func deprecatedResourceSysdigAlertGroupOutlierCreate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := deprecatedGroupOutlierAlertFromResourceData(data)
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

func deprecatedResourceSysdigAlertGroupOutlierUpdate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := deprecatedGroupOutlierAlertFromResourceData(data)
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

func deprecatedResourceSysdigAlertGroupOutlierRead(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
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

	err = deprecatedGroupOutlierAlertToResourceData(&alert, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deprecatedResourceSysdigAlertGroupOutlierDelete(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
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

func deprecatedGroupOutlierAlertFromResourceData(data *schema.ResourceData) (alert *v2.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}

	alert.Type = "HOST_COMPARISON"

	for _, metric := range data.Get("monitor").([]any) {
		alert.Monitor = append(alert.Monitor, &v2.Monitor{
			Metric:       metric.(string),
			StdDevFactor: 2,
		})
	}

	alert.SegmentCondition = &v2.SegmentCondition{Type: "ANY"}
	alert.SegmentBy = []string{"host.mac"}

	return
}

func deprecatedGroupOutlierAlertToResourceData(alert *v2.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	monitorMetrics := []string{}
	for _, v := range alert.Monitor {
		monitorMetrics = append(monitorMetrics, v.Metric)
	}
	_ = data.Set("monitor", monitorMetrics)

	return
}
