package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorAlertGroupOutlier() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigAlertGroupOutlierCreate,
		Update: resourceSysdigAlertGroupOutlierUpdate,
		Read:   resourceSysdigAlertGroupOutlierRead,
		Delete: resourceSysdigAlertGroupOutlierDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createAlertSchema(map[string]*schema.Schema{
			"monitor": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func resourceSysdigAlertGroupOutlierCreate(data *schema.ResourceData, i interface{}) error {
	client := i.(*SysdigClients).sysdigMonitorClient

	alert, err := groupOutlierAlertFromResourceData(data)
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

func resourceSysdigAlertGroupOutlierUpdate(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	alert, err := groupOutlierAlertFromResourceData(data)
	if err != nil {
		return
	}

	alert.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateAlert(*alert)

	return
}

func resourceSysdigAlertGroupOutlierRead(data *schema.ResourceData, i interface{}) (err error) {
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

	err = groupOutlierAlertToResourceData(&alert, data)
	if err != nil {
		return
	}

	return
}

func resourceSysdigAlertGroupOutlierDelete(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return
	}

	return client.DeleteAlert(id)
}

func groupOutlierAlertFromResourceData(data *schema.ResourceData) (alert *monitor.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}

	alert.Type = "HOST_COMPARISON"

	for _, metric := range data.Get("monitor").([]interface{}) {
		alert.Monitor = append(alert.Monitor, &monitor.Monitor{
			Metric:       metric.(string),
			StdDevFactor: 2,
		})
	}

	alert.SegmentCondition = &monitor.SegmentCondition{Type: "ANY"}
	alert.SegmentBy = []string{"host.mac"}

	return
}

func groupOutlierAlertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	monitor_metrics := []string{}
	for _, v := range alert.Monitor {
		monitor_metrics = append(monitor_metrics, v.Metric)
	}
	data.Set("monitor", monitor_metrics)

	return
}
