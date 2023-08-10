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

func resourceSysdigMonitorAlertV2FormBasedPrometheus() *schema.Resource {

	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorAlertV2FormBasedPrometheusCreate,
		UpdateContext: resourceSysdigMonitorAlertV2FormBasedPrometheusUpdate,
		ReadContext:   resourceSysdigMonitorAlertV2FormBasedPrometheusRead,
		DeleteContext: resourceSysdigMonitorAlertV2FormBasedPrometheusDelete,
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
			"query": {
				Type:     schema.TypeString,
				Required: true,
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

func getAlertV2FormBasedPrometheusClient(c SysdigClients) (v2.AlertV2FormBasedPrometheusInterface, error) {
	return getAlertV2Client(c)
}

func resourceSysdigMonitorAlertV2FormBasedPrometheusCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2FormBasedPrometheusClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2FormBasedPrometheusStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	aCreated, err := client.CreateAlertV2FormBasedPrometheus(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(aCreated.ID))

	err = updateAlertV2FormBasedPrometheusState(d, &aCreated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2FormBasedPrometheusRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2FormBasedPrometheusClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := client.GetAlertV2FormBasedPrometheus(ctx, id)
	if err != nil {
		if err == v2.AlertV2NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = updateAlertV2FormBasedPrometheusState(d, &a)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2FormBasedPrometheusUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2FormBasedPrometheusClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2FormBasedPrometheusStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	a.ID, _ = strconv.Atoi(d.Id())

	aUpdated, err := client.UpdateAlertV2FormBasedPrometheus(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateAlertV2FormBasedPrometheusState(d, &aUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2FormBasedPrometheusDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getAlertV2FormBasedPrometheusClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlertV2FormBasedPrometheus(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildAlertV2FormBasedPrometheusStruct(d *schema.ResourceData) (*v2.AlertV2FormBasedPrometheus, error) {
	alertV2Common := buildAlertV2CommonStruct(d)
	alertV2Common.Type = string(v2.AlertV2TypeFormBasedPrometheus)
	config := v2.AlertV2ConfigFormBasedPrometheus{}

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

	//Query
	config.Query = d.Get("query").(string)

	config.NoDataBehaviour = d.Get("no_data_behaviour").(string)

	alert := &v2.AlertV2FormBasedPrometheus{
		AlertV2Common: *alertV2Common,
		DurationSec:   0,
		Config:        config,
	}
	return alert, nil
}

func updateAlertV2FormBasedPrometheusState(d *schema.ResourceData, alert *v2.AlertV2FormBasedPrometheus) error {
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

	_ = d.Set("query", alert.Config.Query)

	_ = d.Set("no_data_behaviour", alert.Config.NoDataBehaviour)

	return nil
}
