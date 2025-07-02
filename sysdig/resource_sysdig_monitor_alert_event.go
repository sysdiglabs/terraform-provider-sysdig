package sysdig

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigMonitorAlertEvent() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		DeprecationMessage: "\"sysdig_monitor_alert_event\" has been deprecated and will be removed in future releases, use \"sysdig_monitor_alert_v2_event\" instead",
		CreateContext:      resourceSysdigAlertEventCreate,
		UpdateContext:      resourceSysdigAlertEventUpdate,
		ReadContext:        resourceSysdigAlertEventRead,
		DeleteContext:      resourceSysdigAlertEventDelete,
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
			"event_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_rel": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"multiple_alerts_by": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func resourceSysdigAlertEventCreate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := eventAlertFromResourceData(data)
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

func resourceSysdigAlertEventUpdate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getMonitorAlertClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := eventAlertFromResourceData(data)
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

func resourceSysdigAlertEventRead(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
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

	err = eventAlertToResourceData(&alert, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigAlertEventDelete(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
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

func eventAlertFromResourceData(data *schema.ResourceData) (alert *v2.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}

	eventRel := data.Get("event_rel").(string)
	eventCount := data.Get("event_count").(int)
	alert.Condition = fmt.Sprintf("count(customEvent) %s %d", eventRel, eventCount)
	alert.Type = "EVENT"
	alert.Criteria = &v2.Criteria{
		Text:   data.Get("event_name").(string),
		Source: data.Get("source").(string),
	}

	if alertsBy, ok := data.GetOk("multiple_alerts_by"); ok {
		alert.SegmentCondition = &v2.SegmentCondition{Type: "ANY"}
		for _, v := range alertsBy.([]any) {
			alert.SegmentBy = append(alert.SegmentBy, v.(string))
		}
	}

	return
}

// https://regex101.com/r/79VIkC/1
var alertConditionRegex = regexp.MustCompile(`count\(customEvent\)\s*(?P<rel>[^\w\s]+)\s*(?P<count>\d+)`)

func eventAlertToResourceData(alert *v2.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	relIndex := alertConditionRegex.SubexpIndex("rel")
	countIndex := alertConditionRegex.SubexpIndex("count")
	matches := alertConditionRegex.FindStringSubmatch(alert.Condition)
	if matches == nil {
		return fmt.Errorf("alert condition %s does not match expected expression %s", alert.Condition, alertConditionRegex.String())
	}

	eventRel := matches[relIndex]
	eventCount, err := strconv.Atoi(matches[countIndex])
	if err != nil {
		return
	}

	_ = data.Set("event_rel", eventRel)
	_ = data.Set("event_count", eventCount)
	_ = data.Set("event_name", alert.Criteria.Text)
	_ = data.Set("source", alert.Criteria.Source)
	_ = data.Set("multiple_alerts_by", alert.SegmentBy)

	return
}
