package sysdig

import (
	"errors"
	"maps"
	"regexp"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		"group_name": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "default",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return strings.EqualFold(old, new)
			},
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
			Type:     schema.TypeSet,
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

	maps.Copy(alertSchema, original)

	return alertSchema
}

func alertFromResourceData(d *schema.ResourceData) (alert *v2.Alert, err error) {
	triggerAfterMinutes := time.Duration(d.Get("trigger_after_minutes").(int)) * time.Minute
	timespan := int(triggerAfterMinutes.Microseconds())
	alert = &v2.Alert{
		Name:                   d.Get("name").(string),
		Type:                   "MANUAL",
		Timespan:               &timespan,
		SegmentBy:              []string{},
		NotificationChannelIds: []int{},
		CustomNotification: &v2.CustomNotification{
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
	if groupName, ok := d.GetOk("group_name"); ok {
		alert.GroupName = strings.ToLower(groupName.(string))
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
		channelSet := channels.(*schema.Set)
		for _, channel := range channelSet.List() {
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

func alertToResourceData(alert *v2.Alert, data *schema.ResourceData) (err error) {
	var triggerAfterMinutes int
	if alert.Timespan != nil {
		triggerAfterMinutes = int((time.Duration(*alert.Timespan) * time.Microsecond).Minutes())
	}

	_ = data.Set("version", alert.Version)
	_ = data.Set("name", alert.Name)
	_ = data.Set("description", alert.Description)
	_ = data.Set("scope", alert.Filter)
	_ = data.Set("trigger_after_minutes", triggerAfterMinutes)
	_ = data.Set("group_name", alert.GroupName)
	_ = data.Set("team", alert.TeamID)
	_ = data.Set("enabled", alert.Enabled)
	_ = data.Set("severity", alert.Severity)

	if len(alert.NotificationChannelIds) > 0 {
		_ = data.Set("notification_channels", alert.NotificationChannelIds)
	}

	if alert.ReNotify {
		_ = data.Set("renotification_minutes", alert.ReNotifyMinutes)
	}

	if alert.CustomNotification != nil &&
		(alert.CustomNotification.TitleTemplate != defaultAlertTitle || alert.CustomNotification.AppendText != "" || alert.CustomNotification.PrependText != "") {
		customNotification := map[string]any{
			"title": alert.CustomNotification.TitleTemplate,
		}

		if alert.CustomNotification.AppendText != "" {
			customNotification["append"] = alert.CustomNotification.AppendText
		}

		if alert.CustomNotification.PrependText != "" {
			customNotification["prepend"] = alert.CustomNotification.PrependText
		}

		_ = data.Set("custom_notification", []any{customNotification})
	}

	if alert.SysdigCapture != nil && alert.SysdigCapture.Enabled {
		capture := map[string]any{
			"filename": alert.SysdigCapture.Name,
			"duration": alert.SysdigCapture.Duration,
		}
		if alert.SysdigCapture.Filters != "" {
			capture["filters"] = alert.SysdigCapture.Filters
		}
		_ = data.Set("capture", []any{capture})
	}

	return
}

func sysdigCaptureFromSet(d *schema.Set) (captures []*v2.SysdigCapture, err error) {
	for _, v := range d.List() {
		m := v.(map[string]any)
		capture := &v2.SysdigCapture{
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

func getMonitorAlertClient(c SysdigClients) (v2.AlertInterface, error) {
	var client v2.AlertInterface
	var err error
	switch c.GetClientType() {
	case IBMMonitor:
		client, err = c.ibmMonitorClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigMonitorClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}
