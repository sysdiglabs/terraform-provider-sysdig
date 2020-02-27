package sysdig

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorAlertEvent() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigAlertEventCreate,
		Update: resourceSysdigAlertEventUpdate,
		Read:   resourceSysdigAlertEventRead,
		Delete: resourceSysdigAlertEventDelete,

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

func resourceSysdigAlertEventCreate(data *schema.ResourceData, i interface{}) error {
	client := i.(*SysdigClients).sysdigMonitorClient

	alert, err := eventAlertFromResourceData(data)
	if err != nil {
		return err
	}

	alertCreated, err := client.CreateAlert(*alert)
	if err != nil {
		return err
	}

	data.SetId(strconv.Itoa(alertCreated.ID))
	data.Set("version", alertCreated.Version)
	return nil
}

func resourceSysdigAlertEventUpdate(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	alert, err := eventAlertFromResourceData(data)
	if err != nil {
		return
	}

	alert.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateAlert(*alert)

	return
}

func resourceSysdigAlertEventRead(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return
	}

	alert, err := client.GetAlertById(id)

	if err != nil {
		data.SetId("")
		return nil
	}

	err = eventAlertToResourceData(&alert, data)
	if err != nil {
		return
	}

	return
}

func resourceSysdigAlertEventDelete(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return
	}

	return client.DeleteAlert(id)
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

func eventAlertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	var event_rel string
	var event_count int
	_, err = fmt.Sscanf(alert.Condition, "count(customEvent) %s %d", &event_rel, &event_count)
	if err != nil {
		return
	}

	data.Set("event_rel", event_rel)
	data.Set("event_count", event_count)
	data.Set("event_name", alert.Criteria.Text)
	data.Set("source", alert.Criteria.Source)
	data.Set("multiple_alerts_by", alert.SegmentBy)

	return
}
