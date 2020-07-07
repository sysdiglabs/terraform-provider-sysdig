package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
)

func resourceSysdigMonitorNotificationChannelWebhook() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigMonitorNotificationChannelWebhookCreate,
		Update: resourceSysdigMonitorNotificationChannelWebhookUpdate,
		Read:   resourceSysdigMonitorNotificationChannelWebhookRead,
		Delete: resourceSysdigMonitorNotificationChannelWebhookDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceSysdigMonitorNotificationChannelWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	notificationChannel, err := monitorNotificationChannelWebhookFromResourceData(d)
	if err != nil {
		return err
	}

	notificationChannel, err = client.CreateNotificationChannel(notificationChannel)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))
	d.Set("version", notificationChannel.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigMonitorNotificationChannelWebhookRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(id)

	if err != nil {
		d.SetId("")
	}

	err = monitorNotificationChannelWebhookToResourceData(&nc, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	nc, err := monitorNotificationChannelWebhookFromResourceData(d)
	if err != nil {
		return err
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(nc)

	return err
}

func resourceSysdigMonitorNotificationChannelWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteNotificationChannel(id)
}

// Channel type for Notification Channels

func monitorNotificationChannelWebhookFromResourceData(d *schema.ResourceData) (nc monitor.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = "WEBHOOK"
	nc.Options.Url = d.Get("url").(string)
	return
}

func monitorNotificationChannelWebhookToResourceData(nc *monitor.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	d.Set("url", nc.Options.Url)
	return
}
