package sysdig

import (
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createSecureNotificationChannelSchema(original map[string]*schema.Schema) map[string]*schema.Schema {
	notificationChannelSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"share_with_current_team": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"notify_when_ok": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"notify_when_resolved": {
			Type:     schema.TypeBool,
			Required: true,
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

	nc = v2.NotificationChannel{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
		TeamID:  tID,
		Options: v2.NotificationChannelOptions{
			NotifyOnOk:           d.Get("notify_when_ok").(bool),
			NotifyOnResolve:      d.Get("notify_when_resolved").(bool),
			SendTestNotification: d.Get("send_test_notification").(bool),
		},
	}
	return
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
	_ = data.Set("notify_when_ok", nc.Options.NotifyOnOk)
	_ = data.Set("notify_when_resolved", nc.Options.NotifyOnResolve)
	_ = data.Set("send_test_notification", nc.Options.SendTestNotification)

	return
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
