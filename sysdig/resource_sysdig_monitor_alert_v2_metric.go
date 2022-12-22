package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
)

func resourceSysdigMonitorAlertV2Metric() *schema.Resource {

	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorAlertV2MetricCreate,
		UpdateContext: resourceSysdigMonitorAlertV2MetricUpdate,
		ReadContext:   resourceSysdigMonitorAlertV2MetricRead,
		DeleteContext: resourceSysdigMonitorAlertV2MetricDelete,
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
			"operator": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{">", ">=", "<", "<=", "=", "!="}, false),
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"warning_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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
		})),
	}
}

func resourceSysdigMonitorAlertV2MetricCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2MetricStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	aCreated, err := client.CreateAlertV2Metric(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(aCreated.ID))

	err = updateAlertV2MetricState(d, &aCreated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2MetricRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := client.GetAlertV2MetricById(ctx, id)

	if err != nil {
		d.SetId("")
		return nil
	}

	err = updateAlertV2MetricState(d, &a)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2MetricUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2MetricStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	a.ID, _ = strconv.Atoi(d.Id())

	aUpdated, err := client.UpdateAlertV2Metric(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateAlertV2MetricState(d, &aUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2MetricDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlertV2Metric(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildAlertV2MetricStruct(d *schema.ResourceData) (*monitor.AlertV2Metric, error) {
	alertV2Common := buildAlertV2CommonStruct(d)
	alertV2Common.Type = monitor.AlertV2AlertType_Manual
	config := monitor.AlertV2ConfigMetric{}

	buildScopedSegmentedConfigStruct(d, &config.ScopedSegmentedConfig)

	//ConditionOperator
	config.ConditionOperator = d.Get("operator").(string)

	//threshold
	config.Threshold = d.Get("threshold").(float64)

	//WarningThreshold
	if warningThreshold, ok := d.GetOk("warning_threshold"); ok {
		wts := warningThreshold.(string)
		wt, err := strconv.ParseFloat(wts, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert warning_threshold to a number: %w", err)
		}
		config.WarningThreshold = &wt
		config.WarningConditionOperator = config.ConditionOperator
	}

	//TimeAggregation
	config.TimeAggregation = d.Get("time_aggregation").(string)

	//GroupAggregation
	config.GroupAggregation = d.Get("group_aggregation").(string)

	//Metric
	metric := d.Get("metric").(string)
	config.Metric.ID = metric

	config.NoDataBehaviour = d.Get("no_data_behaviour").(string)

	alert := &monitor.AlertV2Metric{
		AlertV2Common: *alertV2Common,
		Config:        config,
	}
	return alert, nil
}

func updateAlertV2MetricState(d *schema.ResourceData, alert *monitor.AlertV2Metric) error {
	err := updateAlertV2CommonState(d, &alert.AlertV2Common)
	if err != nil {
		return err
	}

	err = updateScopedSegmentedConfigState(d, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return err
	}

	_ = d.Set("operator", alert.Config.ConditionOperator)

	_ = d.Set("threshold", alert.Config.Threshold)

	if alert.Config.WarningThreshold != nil {
		_ = d.Set("warning_threshold", fmt.Sprintf("%v", *alert.Config.WarningThreshold))
	}

	_ = d.Set("time_aggregation", alert.Config.TimeAggregation)

	_ = d.Set("group_aggregation", alert.Config.GroupAggregation)

	_ = d.Set("metric", alert.Config.Metric.ID)

	_ = d.Set("no_data_behaviour", alert.Config.NoDataBehaviour)

	return nil
}
