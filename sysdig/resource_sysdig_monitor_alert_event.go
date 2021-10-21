package sysdig

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
)

func resourceSysdigMonitorAlertEvent() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigAlertEventCreate,
		UpdateContext: resourceSysdigAlertEventUpdate,
		ReadContext:   resourceSysdigAlertEventRead,
		DeleteContext: resourceSysdigAlertEventDelete,
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

func resourceSysdigAlertEventCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
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

	err = data.Set("version", alertCreated.Version)
	if err != nil {
		log.Println("error asigning 'version' to alert")
	}
	return nil
}

func resourceSysdigAlertEventUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
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

func resourceSysdigAlertEventRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := client.GetAlertById(ctx, id)

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

func resourceSysdigAlertEventDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
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

func eventAlertFromResourceData(data *schema.ResourceData) (alert *monitor.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}

	event_rel := data.Get("event_rel").(string)
	event_count := data.Get("event_count").(int)
	alert.Condition = fmt.Sprintf("count(customEvent) %s %d", event_rel, event_count)
	alert.Type = "EVENT"
	alert.Criteria = &monitor.Criteria{
		Text:   data.Get("event_name").(string),
		Source: data.Get("source").(string),
	}

	if alerts_by, ok := data.GetOk("multiple_alerts_by"); ok {
		alert.SegmentCondition = &monitor.SegmentCondition{Type: "ANY"}
		for _, v := range alerts_by.([]interface{}) {
			alert.SegmentBy = append(alert.SegmentBy, v.(string))
		}
	}

	return
}

// https://regex101.com/r/79VIkC/1
var alertConditionRegex = regexp.MustCompile(`count\(customEvent\)\s*(?P<rel>[^\w\s]+)\s*(?P<count>\d+)`)

func eventAlertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
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

	event_rel := matches[relIndex]
	event_count, err := strconv.Atoi(matches[countIndex])
	if err != nil {
		return
	}

	err = data.Set("event_rel", event_rel)
	if err != nil {
		log.Println("error asigning 'event_rel' to alert")
	}
	err = data.Set("event_count", event_count)
	if err != nil {
		log.Println("error asigning 'event_count' to alert")
	}
	err = data.Set("event_name", alert.Criteria.Text)
	if err != nil {
		log.Println("error asigning 'event_name' to alert")
	}
	err = data.Set("source", alert.Criteria.Source)
	if err != nil {
		log.Println("error asigning 'source' to alert")
	}
	err = data.Set("multiple_alerts_by", alert.SegmentBy)
	if err != nil {
		log.Println("error asigning 'multiple_alerts_by' to alert")
	}

	return
}
