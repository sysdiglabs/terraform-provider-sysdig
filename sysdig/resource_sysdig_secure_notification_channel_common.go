package sysdig

import (
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func createSecureNotificationChannelSchema(original map[string]*schema.Schema) map[string]*schema.Schema {
	notificationChannelSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"share_with_current_team": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"notify_when_ok": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			Deprecated:   "The notify_when_ok field is deprecated and will be removed in a future version. This flag has been replaced by the `notify_on_resolve` field inside `notification_channels` when defining an alert resource.",
		},
		"notify_when_resolved": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			Deprecated:   "The notify_when_resolved field is deprecated and will be removed in a future version. This flag has been replaced by the `notify_on_acknowledge` field inside `notification_channels` when defining an alert resource.",
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"send_test_notification": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}

	for k, v := range original {
		notificationChannelSchema[k] = v
	}

	return notificationChannelSchema
}

func secureNotificationChannelFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	var tID *int
	shareWithCurrentTeam := d.Get("share_with_current_team").(bool)
	if shareWithCurrentTeam {
		tID = &teamID
	}

	var notifyOnOk *bool = nil
	if onOk, ok := d.GetOk("notify_when_ok"); ok && onOk.(string) != "" {
		if onOk.(string) == "true" {
			trueValue := true
			notifyOnOk = &trueValue
		} else {
			falseValue := false
			notifyOnOk = &falseValue
		}
	}

	var notifyOnResolve *bool = nil
	if onResolve, ok := d.GetOk("notify_when_resolved"); ok && onResolve.(string) != "" {
		if onResolve.(string) == "true" {
			trueValue := true
			notifyOnResolve = &trueValue
		} else {
			falseValue := false
			notifyOnResolve = &falseValue
		}
	}

	nc = v2.NotificationChannel{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
		TeamID:  tID,
		Options: v2.NotificationChannelOptions{
			NotifyOnOk:           notifyOnOk,
			NotifyOnResolve:      notifyOnResolve,
			SendTestNotification: d.Get("send_test_notification").(bool),
		},
	}
	return nc, err
}

func secureNotificationChannelToResourceData(nc *v2.NotificationChannel, data *schema.ResourceData) (err error) {
	_ = data.Set("version", nc.Version)
	_ = data.Set("name", nc.Name)
	_ = data.Set("enabled", nc.Enabled)
	var shareWithCurrentTeam bool
	if nc.TeamID != nil {
		shareWithCurrentTeam = true
	}

	err = data.Set("share_with_current_team", shareWithCurrentTeam)
	if err != nil {
		return err
	}

	if nc.Options.NotifyOnOk != nil {
		if *nc.Options.NotifyOnOk {
			_ = data.Set("notify_when_ok", "true")
		} else {
			_ = data.Set("notify_when_ok", "false")
		}
	}

	if nc.Options.NotifyOnResolve != nil {
		if *nc.Options.NotifyOnResolve {
			_ = data.Set("notify_when_resolved", "true")
		} else {
			_ = data.Set("notify_when_resolved", "false")
		}
	}

	// do not update "send_test_notification" from the api response as it will always be "false" on subsequent reads because the fields is not persisted

	return err
}

func getSecureNotificationChannelClient(c SysdigClients) (v2.NotificationChannelInterface, error) {
	var client v2.NotificationChannelInterface
	var err error
	switch c.GetClientType() {
	case IBMSecure:
		client, err = c.ibmSecureClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigSecureClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}
