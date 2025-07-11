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

func resourceSysdigMonitorAlertV2Event() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorAlertV2EventCreate,
		UpdateContext: resourceSysdigMonitorAlertV2EventUpdate,
		ReadContext:   resourceSysdigMonitorAlertV2EventRead,
		DeleteContext: resourceSysdigMonitorAlertV2EventDelete,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true, // computed if range_seconds is defined
				Deprecated:   "Use range_seconds instead",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"range_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true, // computed if trigger_after_minutes is defined
				ExactlyOneOf: []string{"trigger_after_minutes"},
				ValidateFunc: validation.IntAtLeast(60),
			},
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
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		})),
	}
}

func getAlertV2EventClient(c SysdigClients) (v2.AlertV2EventInterface, error) {
	return getAlertV2Client(c)
}

func resourceSysdigMonitorAlertV2EventCreate(ctx context.Context, d *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getAlertV2EventClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2EventStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	aCreated, err := client.CreateAlertV2Event(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(aCreated.ID))

	err = updateAlertV2EventState(d, &aCreated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2EventRead(ctx context.Context, d *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getAlertV2EventClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := client.GetAlertV2EventByID(ctx, id)
	if err != nil {
		if err == v2.ErrAlertV2NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	err = updateAlertV2EventState(d, &a)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2EventUpdate(ctx context.Context, d *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getAlertV2EventClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2EventStruct(d)
	if err != nil {
		return diag.FromErr(err)
	}

	a.ID, _ = strconv.Atoi(d.Id())

	aUpdated, err := client.UpdateAlertV2Event(ctx, *a)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateAlertV2EventState(d, &aUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2EventDelete(ctx context.Context, d *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getAlertV2EventClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlertV2Event(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildAlertV2EventStruct(d *schema.ResourceData) (*v2.AlertV2Event, error) {
	alertV2Common := buildAlertV2CommonStruct(d)
	alertV2Common.Type = string(v2.AlertV2TypeEvent)
	config := v2.AlertV2ConfigEvent{}

	buildScopedSegmentedConfigStruct(d, &config.ScopedSegmentedConfig)

	// ConditionOperator
	config.ConditionOperator = d.Get("operator").(string)

	// threshold
	config.Threshold = d.Get("threshold").(float64)

	// WarningThreshold
	if warningThreshold, ok := d.GetOk("warning_threshold"); ok {
		wts := warningThreshold.(string)
		wt, err := strconv.ParseFloat(wts, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert warning_threshold to a number: %w", err)
		}
		config.WarningThreshold = &wt
		config.WarningConditionOperator = config.ConditionOperator
	}

	// filter
	config.Filter = d.Get("filter").(string)

	// tags
	tags := make([]string, 0)
	if sources, ok := d.GetOk("sources"); ok {
		sourcesList := sources.(*schema.Set).List()
		for _, s := range sourcesList {
			tags = append(tags, s.(string))
		}
	}
	config.Tags = tags

	if attr, ok := d.GetOk("range_seconds"); ok && attr != nil {
		config.Range = d.Get("range_seconds").(int)
	}

	if d.HasChange("trigger_after_minutes") {
		// GetOk returns true even if the value is stored only in the state and not in the user config:
		// to avoid applying a trigger_after_minutes old value from the state even if the user removed it from the config
		// we use HasChange that is true only if the user has changed (or created) it - and so it must be in the config
		if attr, ok := d.GetOk("trigger_after_minutes"); ok && attr != nil {
			config.Range = minutesToSeconds(d.Get("trigger_after_minutes").(int))
		}
	}

	alert := &v2.AlertV2Event{
		AlertV2Common: *alertV2Common,
		Config:        config,
	}
	return alert, nil
}

func updateAlertV2EventState(d *schema.ResourceData, alert *v2.AlertV2Event) error {
	err := updateAlertV2CommonState(d, &alert.AlertV2Common)
	if err != nil {
		return err
	}

	err = updateScopedSegmentedConfigState(d, &alert.Config.ScopedSegmentedConfig)
	if err != nil {
		return err
	}

	_ = d.Set("trigger_after_minutes", secondsToMinutes(alert.Config.Range))
	_ = d.Set("range_seconds", alert.Config.Range)

	_ = d.Set("operator", alert.Config.ConditionOperator)

	_ = d.Set("threshold", alert.Config.Threshold)

	if alert.Config.WarningThreshold != nil {
		_ = d.Set("warning_threshold", fmt.Sprintf("%v", *alert.Config.WarningThreshold))
	}

	_ = d.Set("filter", alert.Config.Filter)

	if len(alert.Config.Tags) > 0 {
		_ = d.Set("sources", alert.Config.Tags)
	}

	return nil
}
