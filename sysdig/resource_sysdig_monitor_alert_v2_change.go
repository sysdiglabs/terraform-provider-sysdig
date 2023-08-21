package sysdig

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func toIntSec(t time.Duration) int {
	return int(t.Seconds())
}

var allowedTimeRanges = map[int]map[int]struct{}{ // for each shorter time range, a set of allowed longer time ranges
	toIntSec(time.Minute * 5): {
		toIntSec(time.Hour):     {},
		toIntSec(time.Hour * 2): {},
		toIntSec(time.Hour * 3): {},
	},
	toIntSec(time.Minute * 10): {
		toIntSec(time.Hour):     {},
		toIntSec(time.Hour * 2): {},
		toIntSec(time.Hour * 3): {},
		toIntSec(time.Hour * 4): {},
		toIntSec(time.Hour * 5): {},
		toIntSec(time.Hour * 6): {},
		toIntSec(time.Hour * 7): {},
		toIntSec(time.Hour * 8): {},
	},
	toIntSec(time.Hour): {
		toIntSec(time.Hour * 4):  {},
		toIntSec(time.Hour * 5):  {},
		toIntSec(time.Hour * 6):  {},
		toIntSec(time.Hour * 7):  {},
		toIntSec(time.Hour * 8):  {},
		toIntSec(time.Hour * 9):  {},
		toIntSec(time.Hour * 10): {},
		toIntSec(time.Hour * 11): {},
		toIntSec(time.Hour * 12): {},
		toIntSec(time.Hour * 13): {},
		toIntSec(time.Hour * 14): {},
		toIntSec(time.Hour * 15): {},
		toIntSec(time.Hour * 16): {},
		toIntSec(time.Hour * 17): {},
		toIntSec(time.Hour * 18): {},
		toIntSec(time.Hour * 19): {},
		toIntSec(time.Hour * 20): {},
		toIntSec(time.Hour * 21): {},
		toIntSec(time.Hour * 22): {},
		toIntSec(time.Hour * 23): {},
		toIntSec(time.Hour * 24): {},
	},
	toIntSec(time.Hour * 4): {
		toIntSec(time.Hour * 24):     {},
		toIntSec(time.Hour * 24 * 2): {},
		toIntSec(time.Hour * 24 * 3): {},
		toIntSec(time.Hour * 24 * 4): {},
		toIntSec(time.Hour * 24 * 5): {},
		toIntSec(time.Hour * 24 * 6): {},
		toIntSec(time.Hour * 24 * 7): {},
	},
	toIntSec(time.Hour * 24): {
		toIntSec(time.Hour * 24 * 7): {},
	},
}

func resourceSysdigMonitorAlertV2Change() *schema.Resource {

	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorAlertV2ChangeCreate,
		UpdateContext: resourceSysdigMonitorAlertV2ChangeUpdate,
		ReadContext:   resourceSysdigMonitorAlertV2ChangeRead,
		DeleteContext: resourceSysdigMonitorAlertV2ChangeDelete,
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
			"shorter_time_range_seconds": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"longer_time_range_seconds": {
				Type:     schema.TypeInt,
				Required: true,
			},
		})),

		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
			shorterTimeRangeSeconds := diff.Get("shorter_time_range_seconds").(int)
			longerTimeRangeSeconds := diff.Get("longer_time_range_seconds").(int)

			if _, ok := allowedTimeRanges[shorterTimeRangeSeconds]; !ok {
				var allowedValues []int
				for k := range allowedTimeRanges {
					allowedValues = append(allowedValues, k)
				}
				sort.Ints(allowedValues)
				return fmt.Errorf("shorter_time_range_seconds can only have one of the following values: %v, provided: %v", allowedValues, shorterTimeRangeSeconds)
			}

			if _, ok := allowedTimeRanges[shorterTimeRangeSeconds][longerTimeRangeSeconds]; !ok {
				var allowedValues []int
				for k := range allowedTimeRanges[shorterTimeRangeSeconds] {
					allowedValues = append(allowedValues, k)
				}
				sort.Ints(allowedValues)
				return fmt.Errorf("longer_time_range_seconds can only have one of the following values if shorter_time_range_seconds is %v: %v, provided: %v", shorterTimeRangeSeconds, allowedValues, longerTimeRangeSeconds)
			}

			return nil
		},
	}
}

func getAlertV2ChangeClient(c SysdigClients) (v2.AlertV2ChangeInterface, error) {
	return getAlertV2Client(c)
}

func resourceSysdigMonitorAlertV2ChangeCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2ChangeClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2ChangeStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	aCreated, err := client.CreateAlertV2Change(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(aCreated.ID))

	err = updateAlertV2ChangeState(d, &aCreated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2ChangeRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2ChangeClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := client.GetAlertV2Change(ctx, id)
	if err != nil {
		if err == v2.AlertV2NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = updateAlertV2ChangeState(d, &a)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2ChangeUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2ChangeClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2ChangeStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	a.ID, _ = strconv.Atoi(d.Id())

	aUpdated, err := client.UpdateAlertV2Change(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateAlertV2ChangeState(d, &aUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2ChangeDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2ChangeClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlertV2Change(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildAlertV2ChangeStruct(d *schema.ResourceData) (*v2.AlertV2Change, error) {
	alertV2Common := buildAlertV2CommonStruct(d)
	alertV2Common.Type = string(v2.AlertV2TypeChange)
	config := v2.AlertV2ConfigChange{}

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

	//ShorterRangeSec
	config.ShorterRangeSec = d.Get("shorter_time_range_seconds").(int)

	//LongerRangeSec
	config.LongerRangeSec = d.Get("longer_time_range_seconds").(int)

	alert := &v2.AlertV2Change{
		AlertV2Common: *alertV2Common,
		DurationSec:   0,
		Config:        config,
	}
	return alert, nil
}

func updateAlertV2ChangeState(d *schema.ResourceData, alert *v2.AlertV2Change) error {
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

	_ = d.Set("shorter_time_range_seconds", alert.Config.ShorterRangeSec)

	_ = d.Set("longer_time_range_seconds", alert.Config.LongerRangeSec)

	return nil
}
