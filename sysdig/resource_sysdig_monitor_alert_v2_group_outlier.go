package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigMonitorAlertV2GroupOutlier() *schema.Resource {
	timeout := 5 * time.Minute

	resource := &schema.Resource{
		CreateContext: resourceSysdigMonitorAlertV2GroupOutlierCreate,
		UpdateContext: resourceSysdigMonitorAlertV2GroupOutlierUpdate,
		ReadContext:   resourceSysdigMonitorAlertV2GroupOutlierRead,
		DeleteContext: resourceSysdigMonitorAlertV2GroupOutlierDelete,
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
			"observation_window_minutes": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"algorithm": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DBSCAN", "MAD"}, false),
			},
			"mad_threshold": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(1, 100),
			},
			"mad_tolerance": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0.5, 10),
			},
			"dbscan_tolerance": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0.5, 10),
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_aggregation": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"avg", "timeAvg", "sum", "min", "max"}, false),
			},
			"group_aggregation": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"avg", "sum", "min", "max"}, false),
			},
			"no_data_behaviour": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DO_NOTHING",
				ValidateFunc: validation.StringInSlice([]string{"DO_NOTHING", "TRIGGER"}, false),
			},
			"unreported_alert_notifications_retention_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(60),
			},
		})),

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
			algorithm := diff.Get("algorithm").(string)
			madThreshold := diff.Get("mad_threshold").(float64)
			madTolerance := diff.Get("mad_tolerance").(float64)
			dbscanTolerance := diff.Get("dbscan_tolerance").(float64)

			if algorithm == "MAD" && madThreshold == 0 && madTolerance == 0 {
				return fmt.Errorf("mad_threshold and mad_tolerance must both be defined and non zero if algorithm = MAD")
			}
			if algorithm == "DBSCAN" && dbscanTolerance == 0 {
				return fmt.Errorf("dbscan_tolerance must be defined and non zero if algorithm = DBSCAN")
			}
			return nil
		},
	}
	// group outlier alert type must be segmented
	resource.Schema["group_by"].Optional = false
	resource.Schema["group_by"].Required = true

	return resource
}

func getAlertV2GroupOutlierClient(c SysdigClients) (v2.AlertV2GroupOutlierInterface, error) {
	return getAlertV2Client(c)
}

func resourceSysdigMonitorAlertV2GroupOutlierCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2GroupOutlierClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2GroupOutlierStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	aCreated, err := client.CreateAlertV2GroupOutlier(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(aCreated.ID))

	err = updateAlertV2GroupOutlierState(d, &aCreated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2GroupOutlierRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2GroupOutlierClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := client.GetAlertV2GroupOutlier(ctx, id)
	if err != nil {
		if err == v2.AlertV2NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = updateAlertV2GroupOutlierState(d, &a)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2GroupOutlierUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2GroupOutlierClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2GroupOutlierStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	a.ID, _ = strconv.Atoi(d.Id())

	aUpdated, err := client.UpdateAlertV2GroupOutlier(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateAlertV2GroupOutlierState(d, &aUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2GroupOutlierDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2GroupOutlierClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlertV2GroupOutlier(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildAlertV2GroupOutlierStruct(d *schema.ResourceData) (*v2.AlertV2GroupOutlier, error) {
	alertV2Common := buildAlertV2CommonStruct(d)
	alertV2Common.Type = string(v2.AlertV2TypeGroupOutlier)
	config := v2.AlertV2ConfigGroupOutlier{}

	buildScopedSegmentedConfigStruct(d, &config.ScopedSegmentedConfig)

	config.Algorithm = d.Get("algorithm").(string)

	config.MadThreshold = d.Get("mad_threshold").(float64)

	config.MadTolerance = d.Get("mad_tolerance").(float64)

	config.DbscanTolerance = d.Get("dbscan_tolerance").(float64)

	metric := d.Get("metric").(string)
	config.Metric.ID = metric

	config.TimeAggregation = d.Get("time_aggregation").(string)

	config.GroupAggregation = d.Get("group_aggregation").(string)

	config.NoDataBehaviour = d.Get("no_data_behaviour").(string)

	config.ObservationWindow = minutesToSeconds(d.Get("observation_window_minutes").(int))

	var unreportedAlertNotificationsRetentionSec *int
	if unreportedAlertNotificationsRetentionSecInterface, ok := d.GetOk("unreported_alert_notifications_retention_seconds"); ok {
		u := unreportedAlertNotificationsRetentionSecInterface.(int)
		unreportedAlertNotificationsRetentionSec = &u
	}

	alert := &v2.AlertV2GroupOutlier{
		AlertV2Common:                            *alertV2Common,
		Config:                                   config,
		UnreportedAlertNotificationsRetentionSec: unreportedAlertNotificationsRetentionSec,
	}
	return alert, nil
}

func updateAlertV2GroupOutlierState(d *schema.ResourceData, alert *v2.AlertV2GroupOutlier) error {
	err := updateAlertV2CommonState(d, &alert.AlertV2Common)
	if err != nil {
		return err
	}

	err = updateScopedSegmentedConfigState(d, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return err
	}

	_ = d.Set("observation_window_minutes", secondsToMinutes(alert.Config.ObservationWindow))

	_ = d.Set("algorithm", alert.Config.Algorithm)

	_ = d.Set("mad_threshold", alert.Config.MadThreshold)

	_ = d.Set("mad_tolerance", alert.Config.MadTolerance)

	_ = d.Set("dbscan_tolerance", alert.Config.DbscanTolerance)

	_ = d.Set("metric", alert.Config.Metric.ID)

	_ = d.Set("time_aggregation", alert.Config.TimeAggregation)

	_ = d.Set("group_aggregation", alert.Config.GroupAggregation)

	_ = d.Set("no_data_behaviour", alert.Config.NoDataBehaviour)

	if alert.UnreportedAlertNotificationsRetentionSec != nil {
		_ = d.Set("unreported_alert_notifications_retention_seconds", *alert.UnreportedAlertNotificationsRetentionSec)
	} else {
		_ = d.Set("unreported_alert_notifications_retention_seconds", nil)
	}

	return nil
}
