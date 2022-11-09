package sysdig

import (
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/monitor"
)

const defaultAlertV2Title = "{{__alert_name__}} is {{__alert_status__}}"

func minutesToSeconds(minutes int) (seconds int) {
	durationMinutes := time.Duration(minutes) * time.Minute
	return int(durationMinutes.Seconds())
}
func secondsToMinutes(seconds int) (minutes int) {
	durationMinutes := time.Duration(seconds) * time.Second
	return int(durationMinutes.Minutes())
}

func createAlertV2Schema(original map[string]*schema.Schema) map[string]*schema.Schema {
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
			Type:         schema.TypeString,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.StringInSlice(monitor.AlertV2Severity_Values(), true),
		},
		"trigger_after_minutes": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"group": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return strings.EqualFold(old, new)
			},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"team": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"notification_channels": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"renotify_every_minutes": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		"custom_notification": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"subject": {
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
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"duration_seconds": {
						Type:     schema.TypeInt,
						Required: false,
						Optional: true,
						Default:  15,
					},
					"storage": {
						Type:     schema.TypeString,
						Optional: false,
						Required: true,
					},
					"filename": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(monitor.AlertV2CaptureFilenameRegexp), "the filename must end in .scap"),
					},
					"filter": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "",
					},
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
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

func buildAlertV2CommonStruct(d *schema.ResourceData) (alert *monitor.AlertV2Common) {

	alert = &monitor.AlertV2Common{
		Name:        d.Get("name").(string),
		Type:        "MANUAL",
		DurationSec: minutesToSeconds(d.Get("trigger_after_minutes").(int)),
		Severity:    d.Get("severity").(string),
		Enabled:     d.Get("enabled").(bool),
	}

	if description, ok := d.GetOk("description"); ok {
		alert.Description = description.(string)
	}

	if group, ok := d.GetOk("group"); ok {
		alert.Group = strings.ToLower(group.(string))
	}

	if version, ok := d.GetOk("version"); ok {
		alert.Version = version.(int)
	}
	if team, ok := d.GetOk("team"); ok {
		alert.TeamID = team.(int)
	}

	alert.NotificationChannelConfigList = &[]monitor.NotificationChannelConfigV2{}
	if attr, ok := d.GetOk("notification_channels"); ok && attr != nil {
		channels := []monitor.NotificationChannelConfigV2{}

		for _, channel := range attr.(*schema.Set).List() {
			channelMap := channel.(map[string]interface{})
			newChannel := monitor.NotificationChannelConfigV2{
				ChannelID: channelMap["id"].(int),
				Type:      channelMap["type"].(string),
			}

			if renotifyEveryMinutes, ok := channelMap["renotify_every_minutes"]; ok {
				newChannel.Options = &monitor.NotificationChannelOptionsV2{
					ReNotifyEverySec: minutesToSeconds(renotifyEveryMinutes.(int)),
				}
			}
			channels = append(channels, newChannel)
		}
		alert.NotificationChannelConfigList = &channels
	}

	if attr, ok := d.GetOk("custom_notification"); ok && attr != nil {
		customNotification := monitor.CustomNotificationTemplateV2{}

		if len(attr.([]interface{})) > 0 {
			m := attr.([]interface{})[0].(map[string]interface{})

			customNotification.Subject = m["subject"].(string)
			customNotification.AppendText = m["append"].(string)
			customNotification.PrependText = m["prepend"].(string)
		} else {
			customNotification.Subject = defaultAlertV2Title
		}

		alert.CustomNotificationTemplate = &customNotification
	}

	if attr, ok := d.GetOk("capture"); ok && attr != nil {
		capture := monitor.CaptureConfigV2{}

		if len(attr.([]interface{})) > 0 {
			m := attr.([]interface{})[0].(map[string]interface{})

			capture.DurationSec = m["duration_seconds"].(int)
			capture.FileName = m["filename"].(string)
			capture.Storage = m["storage"].(string)
			capture.Enabled = m["enabled"].(bool)

			if filter, ok := m["filter"]; ok {
				capture.Filter = filter.(string)
			}
		}
		alert.CaptureConfig = &capture
	}

	return
}

func updateAlertV2CommonState(d *schema.ResourceData, alert *monitor.AlertV2Common) (err error) {
	_ = d.Set("name", alert.Name)
	_ = d.Set("description", alert.Description)
	_ = d.Set("trigger_after_minutes", secondsToMinutes(alert.DurationSec))
	_ = d.Set("severity", alert.Severity)

	// optional with defaults
	_ = d.Set("group", alert.Group)
	_ = d.Set("enabled", alert.Enabled)

	// computed
	_ = d.Set("team", alert.TeamID)
	_ = d.Set("version", alert.Version)

	if alert.NotificationChannelConfigList != nil {
		var notificationChannels []interface{}
		for _, ncc := range *alert.NotificationChannelConfigList {
			config := map[string]interface{}{
				"id":   ncc.ChannelID,
				"type": ncc.Type,
			}

			if ncc.Options != nil {
				config["renotify_every_minutes"] = secondsToMinutes(ncc.Options.ReNotifyEverySec)
			}
			notificationChannels = append(notificationChannels, config)
		}

		_ = d.Set("notification_channels", notificationChannels)
	}

	if alert.CustomNotificationTemplate != nil && !(alert.CustomNotificationTemplate.Subject == defaultAlertV2Title &&
		alert.CustomNotificationTemplate.AppendText == "" &&
		alert.CustomNotificationTemplate.PrependText == "") {
		customNotification := map[string]interface{}{}
		customNotification["subject"] = alert.CustomNotificationTemplate.Subject
		customNotification["append"] = alert.CustomNotificationTemplate.AppendText
		customNotification["prepend"] = alert.CustomNotificationTemplate.PrependText

		_ = d.Set("custom_notification", []interface{}{customNotification})
	}

	if alert.CaptureConfig != nil {
		capture := map[string]interface{}{
			"duration_seconds": alert.CaptureConfig.DurationSec,
			"storage":          alert.CaptureConfig.Storage,
			"filename":         alert.CaptureConfig.FileName,
			"enabled":          alert.CaptureConfig.Enabled,
			"filter":           alert.CaptureConfig.Filter,
		}

		_ = d.Set("capture", []interface{}{capture})
	}

	return nil
}
