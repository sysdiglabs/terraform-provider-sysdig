package sysdig

import (
	"regexp"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const AlertV2CaptureFilenameRegexp = `.*?\.scap`

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
			Default:      string(v2.AlertV2SeverityLow),
			ValidateFunc: validation.StringInSlice(AlertV2SeverityValues(), true),
		},
		"trigger_after_minutes": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"group": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "default",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return strings.EqualFold(old, new)
			},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
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
						Type:       schema.TypeString,
						Optional:   true, //for retro compatibility, content will be discarded, remove this is the next major release
						Default:    "",
						Deprecated: "no need to define \"type\" attribute anymore, please remove it",
					},
					"renotify_every_minutes": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  0,
					},
					"notify_on_resolve": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"main_threshold": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"warning_threshold": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
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
						Optional: true,
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
						Optional: true,
						Default:  15,
					},
					"storage": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "",
					},
					"filename": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(AlertV2CaptureFilenameRegexp), "the filename must end in .scap"), //otherwise the api will silently add .scap at the end
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
		"link": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(AlertLinkV2TypeValues(), true),
					},
					"href": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"id": {
						Type:     schema.TypeString,
						Optional: true,
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

func AlertV2SeverityValues() []string {
	return []string{
		string(v2.AlertV2SeverityHigh),
		string(v2.AlertV2SeverityMedium),
		string(v2.AlertV2SeverityLow),
		string(v2.AlertV2SeverityInfo),
	}
}

func AlertLinkV2TypeValues() []string {
	return []string{
		string(v2.AlertLinkV2TypeDashboard),
		string(v2.AlertLinkV2TypeRunbook),
	}
}

func buildAlertV2CommonStruct(d *schema.ResourceData) *v2.AlertV2Common {

	alert := &v2.AlertV2Common{
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

	alert.NotificationChannelConfigList = []v2.NotificationChannelConfigV2{}
	if attr, ok := d.GetOk("notification_channels"); ok && attr != nil {
		channels := []v2.NotificationChannelConfigV2{}

		for _, channel := range attr.(*schema.Set).List() {
			channelMap := channel.(map[string]interface{})
			newChannel := v2.NotificationChannelConfigV2{
				ChannelID: channelMap["id"].(int),
				//Type: will be added by the sysdig client before the put/post
			}

			if renotifyEveryMinutes, ok := channelMap["renotify_every_minutes"]; ok {
				m := renotifyEveryMinutes.(int)
				if m != 0 {
					s := minutesToSeconds(m)
					newChannel.OverrideOptions.ReNotifyEverySec = &s
				}
			}

			newChannel.OverrideOptions.NotifyOnResolve = channelMap["notify_on_resolve"].(bool)

			newChannel.OverrideOptions.Thresholds = []string{}
			main_threshold := channelMap["main_threshold"].(bool)
			if main_threshold {
				newChannel.OverrideOptions.Thresholds = append(newChannel.OverrideOptions.Thresholds, "MAIN")
			}
			warning_threshold := channelMap["warning_threshold"].(bool)
			if warning_threshold {
				newChannel.OverrideOptions.Thresholds = append(newChannel.OverrideOptions.Thresholds, "WARNING")
			}

			channels = append(channels, newChannel)
		}
		alert.NotificationChannelConfigList = channels
	}

	customNotification := v2.CustomNotificationTemplateV2{}
	if attr, ok := d.GetOk("custom_notification"); ok && attr != nil {
		if len(attr.([]interface{})) > 0 {
			m := attr.([]interface{})[0].(map[string]interface{})

			customNotification.Subject = m["subject"].(string)
			customNotification.AppendText = m["append"].(string)
			customNotification.PrependText = m["prepend"].(string)
		}
	}
	alert.CustomNotificationTemplate = &customNotification

	if attr, ok := d.GetOk("capture"); ok && attr != nil {
		capture := v2.CaptureConfigV2{}

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

	alert.Links = []v2.AlertLinkV2{}
	if attr, ok := d.GetOk("link"); ok && attr != nil {
		for _, link := range attr.(*schema.Set).List() {
			linkMap := link.(map[string]interface{})
			alert.Links = append(alert.Links, v2.AlertLinkV2{
				Type: linkMap["type"].(string),
				Href: linkMap["href"].(string),
				ID:   linkMap["id"].(string), //TODO(dbonf) if referencing a non existing dashboard, API will silently fail (status code: 200) not saving the link, add validation?
			})
		}
	}

	return alert
}

func updateAlertV2CommonState(d *schema.ResourceData, alert *v2.AlertV2Common) (err error) {
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

	var notificationChannels []interface{}
	for _, ncc := range alert.NotificationChannelConfigList {
		config := map[string]interface{}{
			"id":                ncc.ChannelID,
			"notify_on_resolve": ncc.OverrideOptions.NotifyOnResolve,
		}

		if ncc.OverrideOptions.ReNotifyEverySec != nil {
			config["renotify_every_minutes"] = secondsToMinutes(*ncc.OverrideOptions.ReNotifyEverySec)
		} else {
			config["renotify_every_minutes"] = 0
		}

		if ncc.OverrideOptions.Thresholds != nil {
			config["main_threshold"] = false
			config["warning_threshold"] = false
			for _, t := range ncc.OverrideOptions.Thresholds {
				if t == "MAIN" {
					config["main_threshold"] = true
				}
				if t == "WARNING" {
					config["warning_threshold"] = true
				}
			}
		} else {
			// defaults
			config["main_threshold"] = true
			config["warning_threshold"] = false
		}

		notificationChannels = append(notificationChannels, config)
	}
	_ = d.Set("notification_channels", notificationChannels)

	if alert.CustomNotificationTemplate != nil && !(alert.CustomNotificationTemplate.Subject == "" &&
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

	if alert.Links != nil {
		var links []interface{}
		for _, link := range alert.Links {
			links = append(links, map[string]interface{}{
				"type": link.Type,
				"href": link.Href,
				"id":   link.ID,
			})
		}
		_ = d.Set("link", links)
	}

	return nil
}

func createScopedSegmentedAlertV2Schema(original map[string]*schema.Schema) map[string]*schema.Schema {
	sysdigAlertSchema := map[string]*schema.Schema{
		"scope": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"label": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringDoesNotContainAny("."),
					},
					"operator": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"equals", "notEquals", "in", "notIn", "contains", "notContains", "startsWith"}, false),
					},
					"values": {
						Type:     schema.TypeList,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"group_by": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		}}

	for k, v := range original {
		sysdigAlertSchema[k] = v
	}

	return sysdigAlertSchema
}

func buildScopedSegmentedConfigStruct(d *schema.ResourceData, config *v2.ScopedSegmentedConfig) {
	//scope
	expressions := make([]v2.ScopeExpressionV2, 0)
	for _, scope := range d.Get("scope").(*schema.Set).List() {
		scopeMap := scope.(map[string]interface{})
		operator := scopeMap["operator"].(string)
		operand := scopeMap["label"].(string)
		value := make([]string, 0)
		for _, v := range scopeMap["values"].([]interface{}) {
			value = append(value, v.(string))
		}
		expressions = append(expressions, v2.ScopeExpressionV2{
			Operand:  operand, //the sysdig client will rewrite this to be in dot notation
			Operator: operator,
			Value:    value,
		})
	}
	if len(expressions) > 0 {
		config.Scope = &v2.AlertScopeV2{
			Expressions: expressions,
		}
	}

	//SegmentBy
	config.SegmentBy = make([]v2.AlertLabelDescriptorV2, 0)
	labels, ok := d.GetOk("group_by")
	if ok {
		for _, l := range labels.([]interface{}) {
			config.SegmentBy = append(config.SegmentBy, v2.AlertLabelDescriptorV2{
				ID: l.(string), //the sysdig client will rewrite this to be in dot notation
			})
		}
	}
}

func updateScopedSegmentedConfigState(d *schema.ResourceData, config *v2.ScopedSegmentedConfig) error {
	if config.Scope != nil && len(config.Scope.Expressions) > 0 {
		var scope []interface{}
		for _, e := range config.Scope.Expressions {
			// operand possibly holds the old dot notation, we want "label" to be in public notation
			// if the label does not yet exist the descriptor will be empty, use what's in the operand
			label := e.Operand
			if e.Descriptor != nil && e.Descriptor.PublicID != "" {
				label = e.Descriptor.PublicID
			}
			config := map[string]interface{}{
				"label":    label,
				"operator": e.Operator,
				"values":   e.Value,
			}
			scope = append(scope, config)
		}
		_ = d.Set("scope", scope)
	}

	if len(config.SegmentBy) > 0 {
		groups := make([]string, 0)
		for _, s := range config.SegmentBy {
			groups = append(groups, s.PublicID)
		}
		_ = d.Set("group_by", groups)
	}

	return nil
}

func getAlertV2Client(c SysdigClients) (v2.AlertV2Interface, error) {
	var client v2.AlertV2Interface
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
