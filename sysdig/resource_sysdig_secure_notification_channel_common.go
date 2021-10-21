package sysdig

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
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

func secureNotificationChannelFromResourceData(d *schema.ResourceData) (nc secure.NotificationChannel, err error) {
	nc = secure.NotificationChannel{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
		Options: secure.NotificationChannelOptions{
			NotifyOnOk:           d.Get("notify_when_ok").(bool),
			NotifyOnResolve:      d.Get("notify_when_resolved").(bool),
			SendTestNotification: d.Get("send_test_notification").(bool),
		},
	}
	return
}

func secureNotificationChannelToResourceData(nc *secure.NotificationChannel, data *schema.ResourceData) (err error) {
	err = data.Set("version", nc.Version)
	if err != nil {
		log.Println("error assigning 'version'")
	}

	err = data.Set("name", nc.Name)
	if err != nil {
		log.Println("error assigning 'name'")
	}

	err = data.Set("enabled", nc.Enabled)
	if err != nil {
		log.Println("error assigning 'enabled'")
	}

	err = data.Set("notify_when_ok", nc.Options.NotifyOnOk)
	if err != nil {
		log.Println("error assigning 'notify_when_ok'")
	}

	err = data.Set("notify_when_resolved", nc.Options.NotifyOnResolve)
	if err != nil {
		log.Println("error assigning 'notify_when_resolved'")
	}

	err = data.Set("send_test_notification", nc.Options.SendTestNotification)
	if err != nil {
		log.Println("error assigning 'send_test_notification'")
	}
	return
}
