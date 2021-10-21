package sysdig

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
)

const defaultAlertTitle = "{{__alert_name__}} is {{__alert_status__}}"

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
		"custom_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"title": {
						Type:     schema.TypeString,
						Required: true,
					},
					"prepend": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"append": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"capture": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"filename": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`.*?\.scap`), "the filename must end in .scap"),
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
		},
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
			TitleTemplate:  defaultAlertTitle,
			UseNewTemplate: true,
		},
	}

	if _, ok := d.GetOk("custom_notification"); ok {
		if title, ok := d.GetOk("custom_notification.0.title"); ok {
			alert.CustomNotification.TitleTemplate = title.(string)
		}
		if prependText, ok := d.GetOk("custom_notification.0.prepend"); ok {
			alert.CustomNotification.PrependText = prependText.(string)
		}
		if appendText, ok := d.GetOk("custom_notification.0.append"); ok {
			alert.CustomNotification.AppendText = appendText.(string)
		}
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
		for _, channel := range channels.([]interface{}) {
			alert.NotificationChannelIds = append(alert.NotificationChannelIds, channel.(int))
		}
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

	// note: didn't want to change current behaviour, error handing is set to avoid lint errors
	err = data.Set("version", alert.Version)
	if err != nil {
		log.Println("error assigning 'version'")
	}
	err = data.Set("name", alert.Name)
	if err != nil {
		log.Println("error assigning 'name'")
	}
	err = data.Set("description", alert.Description)
	if err != nil {
		log.Println("error assigning 'description'")
	}
	err = data.Set("scope", alert.Filter)
	if err != nil {
		log.Println("error assigning 'scope'")
	}
	err = data.Set("trigger_after_minutes", int(trigger_after_minutes.Minutes()))
	if err != nil {
		log.Println("error assigning 'trigger_after_minutes'")
	}
	err = data.Set("team", alert.TeamID)
	if err != nil {
		log.Println("error assigning 'team'")
	}
	err = data.Set("enabled", alert.Enabled)
	if err != nil {
		log.Println("error assigning 'enabled'")
	}
	err = data.Set("severity", alert.Severity)
	if err != nil {
		log.Println("error assigning 'severity'")
	}

	if len(alert.NotificationChannelIds) > 0 {
		err = data.Set("notification_channels", alert.NotificationChannelIds)
		if err != nil {
			log.Println("error assigning 'notification_channels'")
		}
	}

	if alert.ReNotify {
		err = data.Set("renotification_minutes", alert.ReNotifyMinutes)
		if err != nil {
			log.Println("error assigning 'renotification_minutes'")
		}
	}

	if alert.CustomNotification != nil &&
		(alert.CustomNotification.TitleTemplate != defaultAlertTitle || alert.CustomNotification.AppendText != "" || alert.CustomNotification.PrependText != "") {
		customNotification := map[string]interface{}{
			"title": alert.CustomNotification.TitleTemplate,
		}

		if alert.CustomNotification.AppendText != "" {
			customNotification["append"] = alert.CustomNotification.AppendText
		}

		if alert.CustomNotification.PrependText != "" {
			customNotification["prepend"] = alert.CustomNotification.PrependText
		}

		err = data.Set("custom_notification", []interface{}{customNotification})
		if err != nil {
			log.Println("error assigning 'custom_notification'")
		}
	}

	if alert.SysdigCapture != nil && alert.SysdigCapture.Enabled {
		capture := map[string]interface{}{
			"filename": alert.SysdigCapture.Name,
			"duration": alert.SysdigCapture.Duration,
		}
		if alert.SysdigCapture.Filters != "" {
			capture["filters"] = alert.SysdigCapture.Filters
		}
		err = data.Set("capture", []interface{}{capture})
		if err != nil {
			log.Println("error assigning 'capture'")
		}
	}

	return
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
