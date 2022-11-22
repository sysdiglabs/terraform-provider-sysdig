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

		Schema: createAlertV2Schema(map[string]*schema.Schema{
			"scope": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringDoesNotContainAny("."),
						},
						"op": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"equals", "notEquals", "in", "notIn", "contains", "notContains", "startsWith"}, false),
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"group_by": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"op": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{">", ">=", "<", "<=", "=", "!="}, false),
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"warning_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
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
		}),
	}
}

func resourceSysdigMonitorAlertV2EventCreate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2EventStruct(ctx, d, client)
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

func resourceSysdigMonitorAlertV2EventRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := client.GetAlertV2EventById(ctx, id)

	if err != nil {
		d.SetId("")
		return nil
	}

	err = updateAlertV2EventState(d, &a)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorAlertV2EventUpdate(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	a, err := buildAlertV2EventStruct(ctx, d, client)
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

func resourceSysdigMonitorAlertV2EventDelete(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
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

func buildAlertV2EventStruct(ctx context.Context, d *schema.ResourceData, client monitor.SysdigMonitorClient) (*monitor.AlertV2Event, error) {
	alertV2Common, err := buildAlertV2CommonStruct(ctx, d, client)
	if err != nil {
		return nil, err
	}
	alertV2Common.Type = monitor.AlertV2AlertType_Event
	config := &monitor.AlertV2ConfigEvent{}

	//scope
	expressions := make([]monitor.ScopeExpressionV2, 0)
	for _, scope := range d.Get("scope").(*schema.Set).List() {
		scopeMap := scope.(map[string]interface{})
		operator := scopeMap["op"].(string)
		value := make([]string, 0)
		for _, v := range scopeMap["values"].([]interface{}) {
			value = append(value, v.(string))
		}
		label := scopeMap["label"].(string)
		labelDescriptorV3, err := client.GetLabelDescriptor(ctx, label)
		if err != nil {
			return nil, fmt.Errorf("error getting descriptor for label %s: %w", label, err)
		}
		operand := labelDescriptorV3.LabelDescriptor.ID
		expressions = append(expressions, monitor.ScopeExpressionV2{
			Operand:  operand,
			Operator: operator,
			Value:    value,
		})
	}
	if len(expressions) > 0 {
		config.Scope = &monitor.AlertScopeV2{
			Expressions: expressions,
		}
	}

	//SegmentBy
	config.SegmentBy = make([]monitor.AlertLabelDescriptorV2, 0)
	labels, ok := d.GetOk("group_by")
	if ok {
		for _, l := range labels.([]interface{}) {
			label := l.(string)
			labelDescriptorV3, err := client.GetLabelDescriptor(ctx, label)
			if err != nil {
				return nil, fmt.Errorf("error getting descriptor for label %s: %w", label, err)
			}
			config.SegmentBy = append(config.SegmentBy, monitor.AlertLabelDescriptorV2{
				ID:       labelDescriptorV3.LabelDescriptor.ID,
				PublicID: labelDescriptorV3.LabelDescriptor.PublicID,
			})
		}
	}

	//ConditionOperator
	config.ConditionOperator = d.Get("op").(string)

	//threshold
	config.Threshold = d.Get("threshold").(float64)

	//WarningThreshold
	if warningThreshold, ok := d.GetOk("warning_threshold"); ok {
		wt := warningThreshold.(float64)
		config.WarningThreshold = &wt
		config.WarningConditionOperator = config.ConditionOperator
	}

	//filter
	config.Filter = d.Get("filter").(string)

	//tags
	tags := make([]string, 0)
	if sources, ok := d.GetOk("sources"); ok {
		sourcesList := sources.(*schema.Set).List()
		for _, s := range sourcesList {
			tags = append(tags, s.(string))
		}
	}
	config.Tags = tags

	alert := &monitor.AlertV2Event{
		AlertV2Common: *alertV2Common,
		Config:        config,
	}
	return alert, nil
}

func updateAlertV2EventState(d *schema.ResourceData, alert *monitor.AlertV2Event) error {
	err := updateAlertV2CommonState(d, &alert.AlertV2Common)
	if err != nil {
		return err
	}

	if alert.Config.Scope != nil && len(alert.Config.Scope.Expressions) > 0 {
		var scope []interface{}
		for _, e := range alert.Config.Scope.Expressions {
			// operand possibly holds the old dot notation, we want "label" to be in public notation
			// if the label does not yet exist the descriptor will be empty, use what's in the operand
			label := e.Descriptor.PublicID
			if label == "" {
				label = e.Operand
			}
			config := map[string]interface{}{
				"label":  label,
				"op":     e.Operator,
				"values": e.Value,
			}
			scope = append(scope, config)
		}
		_ = d.Set("scope", scope)
	}

	if len(alert.Config.SegmentBy) > 0 {
		groups := make([]string, 0)
		for _, s := range alert.Config.SegmentBy {
			groups = append(groups, s.PublicID)
		}
		_ = d.Set("group_by", groups)
	}

	_ = d.Set("op", alert.Config.ConditionOperator)

	_ = d.Set("threshold", alert.Config.Threshold)

	if alert.Config.WarningThreshold != nil {
		_ = d.Set("warning_threshold", alert.Config.WarningThreshold)
	}

	_ = d.Set("filter", alert.Config.Filter)

	if len(alert.Config.Tags) > 0 {
		_ = d.Set("sources", alert.Config.Tags)
	}

	return nil
}
