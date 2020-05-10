package sysdig

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorAlertDowntime() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigAlertDowntimeCreate,
		Update: resourceSysdigAlertDowntimeUpdate,
		Read:   resourceSysdigAlertDowntimeRead,
		Delete: resourceSysdigAlertDowntimeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createAlertSchema(map[string]*schema.Schema{
			"entities_to_monitor": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"trigger_after_pct": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
		}),
	}
}

func resourceSysdigAlertDowntimeCreate(data *schema.ResourceData, i interface{}) error {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	alert, err := downtimeAlertFromResourceData(data)
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

func resourceSysdigAlertDowntimeUpdate(data *schema.ResourceData, i interface{}) (err error) {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return
	}

	alert, err := downtimeAlertFromResourceData(data)
	if err != nil {
		return
	}

	alert.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateAlert(*alert)

	return
}

func resourceSysdigAlertDowntimeRead(data *schema.ResourceData, i interface{}) (err error) {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return
	}

	alert, err := client.GetAlertById(id)

	if err != nil {
		data.SetId("")
		return nil
	}

	err = downtimeAlertToResourceData(&alert, data)
	if err != nil {
		return
	}

	return
}
func resourceSysdigAlertDowntimeDelete(data *schema.ResourceData, i interface{}) (err error) {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return
	}

	return client.DeleteAlert(id)
}

func downtimeAlertFromResourceData(d *schema.ResourceData) (alert *monitor.Alert, err error) {
	alert, err = alertFromResourceData(d)
	if err != nil {
		return
	}

	alert.SegmentCondition = &monitor.SegmentCondition{Type: "ANY"}
	alert.Condition = fmt.Sprintf("avg(timeAvg(uptime)) <= %.2f", float64(1.0-(d.Get("trigger_after_pct").(int)/100)))

	entitiesRaw := d.Get("entities_to_monitor").([]interface{})
	for _, entityRaw := range entitiesRaw {
		alert.SegmentBy = append(alert.SegmentBy, entityRaw.(string))
	}

	return
}

func downtimeAlertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	var trigger_after_pct float64
	fmt.Sscanf(alert.Condition, "avg(timeAvg(uptime)) <= %.2f", &trigger_after_pct)
	trigger_after_pct = (1 - trigger_after_pct) * 100

	data.Set("trigger_after_pct", int(trigger_after_pct))
	data.Set("entities_to_monitor", alert.SegmentBy)

	return
}
