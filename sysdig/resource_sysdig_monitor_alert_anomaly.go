package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorAlertAnomaly() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigAlertAnomalyCreate,
		Update: resourceSysdigAlertAnomalyUpdate,
		Read:   resourceSysdigAlertAnomalyRead,
		Delete: resourceSysdigAlertAnomalyDelete,

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
			"multiple_alerts_by": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		}),
	}
}

func resourceSysdigAlertAnomalyCreate(data *schema.ResourceData, i interface{}) error {
	client := i.(*SysdigClients).sysdigMonitorClient

	alert, err := anomalyAlertFromResourceData(data)
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

func resourceSysdigAlertAnomalyUpdate(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	alert, err := anomalyAlertFromResourceData(data)
	if err != nil {
		return
	}

	alert.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateAlert(*alert)

	return
}

func resourceSysdigAlertAnomalyRead(data *schema.ResourceData, i interface{}) (err error) {
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

	err = anomalyAlertToResourceData(&alert, data)
	if err != nil {
		return
	}

	return
}

func resourceSysdigAlertAnomalyDelete(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return
	}

	return client.DeleteAlert(id)
}

func anomalyAlertFromResourceData(data *schema.ResourceData) (alert *monitor.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}

	alert.Type = "BASELINE"

	for _, metric := range data.Get("monitor").([]interface{}) {
		alert.Monitor = append(alert.Monitor, &monitor.Monitor{
			Metric:       metric.(string),
			StdDevFactor: 2,
		})
	}

	if alerts_by, ok := data.GetOk("multiple_alerts_by"); ok {
		alert.SegmentCondition = &monitor.SegmentCondition{Type: "ANY"}
		for _, v := range alerts_by.([]interface{}) {
			alert.SegmentBy = append(alert.SegmentBy, v.(string))
		}
	}

	return
}

func anomalyAlertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	data.Set("multiple_alerts_by", alert.SegmentBy)

	monitor_metrics := []string{}
	for _, v := range alert.Monitor {
		monitor_metrics = append(monitor_metrics, v.Metric)
	}
	data.Set("monitor", monitor_metrics)

	return
}
