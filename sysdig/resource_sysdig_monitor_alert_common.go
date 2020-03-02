package sysdig

import (
	"errors"
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
	"time"
)

func createAlertSchema(original map[string]*schema.Schema) map[string]*schema.Schema {
	alertSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"severity": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(0, 7),
		},
		"trigger_after_minutes": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"scope": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"team": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"notification_channels": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeInt},
			Optional: true,
		},
		"renotification_minutes": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"capture": resourceSysdigMonitorAlertCapture(),
	}

	for k, v := range original {
		alertSchema[k] = v
	}

	return alertSchema
}

func alertFromResourceData(d *schema.ResourceData) (alert *monitor.Alert, err error) {
	trigger_after_minutes := time.Duration(d.Get("trigger_after_minutes").(int)) * time.Minute
	alert = &monitor.Alert{
		Name:                   d.Get("name").(string),
		Type:                   "MANUAL",
		Timespan:               int(trigger_after_minutes.Microseconds()),
		SegmentBy:              []string{},
		NotificationChannelIds: []int{},
		CustomNotification: &monitor.CustomNotification{
			TitleTemplate:  "{{__alert_name__}} is {{__alert_status__}}",
			UseNewTemplate: true,
		},
	}

	if scope, ok := d.GetOk("scope"); ok {
		alert.Filter = scope.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		alert.Description = description.(string)
	}
	if version, ok := d.GetOk("version"); ok {
		alert.Version = version.(int)
	}
	if team, ok := d.GetOk("team"); ok {
		alert.TeamID = team.(int)
	}
	if enabled, ok := d.GetOk("enabled"); ok {
		alert.Enabled = enabled.(bool)
	}

	if channels, ok := d.GetOk("notification_channels"); ok {
		alert.NotificationChannelIds = channels.([]int)
	}

	if renotificationMinutes, ok := d.GetOk("renotification_minutes"); ok {
		alert.ReNotify = true
		alert.ReNotifyMinutes = renotificationMinutes.(int)
	} else {
		alert.ReNotify = false
		alert.ReNotifyMinutes = 1 // Required by the API to be higher than 0 even if the re notification is not set
	}

	if set, ok := d.GetOk("capture"); ok {
		captures, err := sysdigCaptureFromSet(set.(*schema.Set))
		if err != nil {
			return nil, err
		}
		if len(captures) == 0 {
			err = errors.New("capture set is empty")
			return nil, err
		}
		alert.SysdigCapture = captures[0]
	}

	alert.Severity = d.Get("severity").(int)

	return
}

func alertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
	trigger_after_minutes := time.Duration(alert.Timespan) * time.Microsecond

	data.Set("version", alert.Version)
	data.Set("name", alert.Name)
	data.Set("description", alert.Description)
	data.Set("scope", alert.Filter)
	data.Set("trigger_after_minutes", int(trigger_after_minutes.Minutes()))
	data.Set("team", alert.TeamID)
	data.Set("enabled", alert.Enabled)

	if len(alert.NotificationChannelIds) > 0 {
		data.Set("notification_channels", alert.NotificationChannelIds)
	}

	if alert.ReNotify {
		data.Set("renotification_minutes", alert.ReNotifyMinutes)
	}

	if alert.SysdigCapture != nil && alert.SysdigCapture.Enabled {
		capture := map[string]interface{}{
			"filename": alert.SysdigCapture.Name,
			"duration": alert.SysdigCapture.Duration,
		}
		if alert.SysdigCapture.Filters != "" {
			capture["filters"] = alert.SysdigCapture.Filters
		}
		data.Set("capture", []interface{}{capture})
	}

	return
}

func resourceSysdigMonitorAlertCapture() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"filename": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(".*?\\.scap"), "the filename must end in .scap"),
				},
				"duration": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"filter": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
			},
		},
	}
}

func sysdigCaptureFromSet(d *schema.Set) (captures []*monitor.SysdigCapture, err error) {
	for _, v := range d.List() {
		m := v.(map[string]interface{})
		capture := &monitor.SysdigCapture{
			Name:     m["filename"].(string),
			Duration: m["duration"].(int),
			Enabled:  true,
		}
		if filter, ok := m["filter"]; ok {
			capture.Filters = filter.(string)
		}
		captures = append(captures, capture)
	}

	return
}
