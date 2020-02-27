package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorAlertMetric() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigAlertMetricCreate,
		Update: resourceSysdigAlertMetricUpdate,
		Read:   resourceSysdigAlertMetricRead,
		Delete: resourceSysdigAlertMetricDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createAlertSchema(map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeString,
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

func resourceSysdigAlertMetricCreate(data *schema.ResourceData, i interface{}) error {
	client := i.(*SysdigClients).sysdigMonitorClient

	alert, err := metricAlertFromResourceData(data)
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

func resourceSysdigAlertMetricUpdate(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	alert, err := metricAlertFromResourceData(data)
	if err != nil {
		return
	}

	alert.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateAlert(*alert)

	return
}

func resourceSysdigAlertMetricRead(data *schema.ResourceData, i interface{}) (err error) {
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

	err = metricAlertToResourceData(&alert, data)
	if err != nil {
		return
	}

	return
}

func resourceSysdigAlertMetricDelete(data *schema.ResourceData, i interface{}) (err error) {
	client := i.(*SysdigClients).sysdigMonitorClient

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return
	}

	return client.DeleteAlert(id)
}

func metricAlertFromResourceData(data *schema.ResourceData) (alert *monitor.Alert, err error) {
	alert, err = alertFromResourceData(data)
	if err != nil {
		return
	}
	alert.Condition = data.Get("metric").(string)

	if alerts_by, ok := data.GetOk("multiple_alerts_by"); ok {
		alert.SegmentCondition = &monitor.SegmentCondition{Type: "ANY"}
		for _, v := range alerts_by.([]interface{}) {
			alert.SegmentBy = append(alert.SegmentBy, v.(string))
		}
	}
	return
}

func metricAlertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	data.Set("metric", alert.Condition)
	data.Set("multiple_alerts_by", alert.SegmentBy)

	return
}
