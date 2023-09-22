package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigMonitorAlertV2Downtime() *schema.Resource {

	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorAlertV2DowntimeCreate,
		UpdateContext: resourceSysdigMonitorAlertV2DowntimeUpdate,
		ReadContext:   resourceSysdigMonitorAlertV2DowntimeRead,
		DeleteContext: resourceSysdigMonitorAlertV2DowntimeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createScopedSegmentedAlertV2Schema(createAlertV2Schema(map[string]*schema.Schema{
			"trigger_after_minutes": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  100,
			},
			"metric": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"sysdig_container_up", "sysdig_program_up", "sysdig_host_up"}, true),
			},
			"unreported_alert_notifications_retention_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(60),
			},
		})),
	}
}

func getAlertV2DowntimeClient(c SysdigClients) (v2.AlertV2DowntimeInterface, error) {
	return getAlertV2Client(c)
}

func resourceSysdigMonitorAlertV2DowntimeCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2DowntimeClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a := buildAlertV2DowntimeStruct(d)

	aCreated, err := client.CreateAlertV2Downtime(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(aCreated.ID))

	err = updateAlertV2DowntimeState(d, &aCreated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2DowntimeRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2DowntimeClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := client.GetAlertV2Downtime(ctx, id)
	if err != nil {
		if err == v2.AlertV2NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = updateAlertV2DowntimeState(d, &a)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2DowntimeUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2DowntimeClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a := buildAlertV2DowntimeStruct(d)

	a.ID, _ = strconv.Atoi(d.Id())

	aUpdated, err := client.UpdateAlertV2Downtime(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateAlertV2DowntimeState(d, &aUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2DowntimeDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2DowntimeClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlertV2Downtime(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildAlertV2DowntimeStruct(d *schema.ResourceData) *v2.AlertV2Downtime {
	alertV2Common := buildAlertV2CommonStruct(d)
	alertV2Common.Type = string(v2.AlertV2TypeManual)
	config := v2.AlertV2ConfigDowntime{}

	buildScopedSegmentedConfigStruct(d, &config.ScopedSegmentedConfig)

	//TimeAggregation
	config.TimeAggregation = "timeAvg"

	//GroupAggregation
	config.GroupAggregation = "avg"

	//ConditionOperator
	config.ConditionOperator = "<="

	//threshold
	config.Threshold = 1 - d.Get("threshold").(float64)/100

	//Downtime
	metric := d.Get("metric").(string)
	config.Metric.ID = metric

	config.NoDataBehaviour = "DO_NOTHING"

	var unreportedAlertNotificationsRetentionSec *int
	if unreportedAlertNotificationsRetentionSecInterface, ok := d.GetOk("unreported_alert_notifications_retention_seconds"); ok {
		u := unreportedAlertNotificationsRetentionSecInterface.(int)
		unreportedAlertNotificationsRetentionSec = &u
	}

	alert := &v2.AlertV2Downtime{
		AlertV2Common:                            *alertV2Common,
		DurationSec:                              minutesToSeconds(d.Get("trigger_after_minutes").(int)),
		Config:                                   config,
		UnreportedAlertNotificationsRetentionSec: unreportedAlertNotificationsRetentionSec,
	}
	return alert
}

func updateAlertV2DowntimeState(d *schema.ResourceData, alert *v2.AlertV2Downtime) error {
	err := updateAlertV2CommonState(d, &alert.AlertV2Common)
	if err != nil {
		return err
	}

	err = updateScopedSegmentedConfigState(d, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return err
	}

	_ = d.Set("trigger_after_minutes", secondsToMinutes(alert.DurationSec))

	_ = d.Set("threshold", (1-alert.Config.Threshold)*100)

	_ = d.Set("metric", alert.Config.Metric.ID)

	if alert.UnreportedAlertNotificationsRetentionSec != nil {
		_ = d.Set("unreported_alert_notifications_retention_seconds", *alert.UnreportedAlertNotificationsRetentionSec)
	} else {
		_ = d.Set("unreported_alert_notifications_retention_seconds", nil)
	}

	return nil
}
