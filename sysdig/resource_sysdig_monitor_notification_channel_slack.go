package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
)

func resourceSysdigMonitorNotificationChannelSlack() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigMonitorNotificationChannelSlackCreate,
		Update: resourceSysdigMonitorNotificationChannelSlackUpdate,
		Read:   resourceSysdigMonitorNotificationChannelSlackRead,
		Delete: resourceSysdigMonitorNotificationChannelSlackDelete,

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
			"channel": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceSysdigMonitorNotificationChannelSlackCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	notificationChannel, err := monitorNotificationChannelSlackFromResourceData(d)
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
func resourceSysdigMonitorNotificationChannelSlackRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(id)

	if err != nil {
		d.SetId("")
	}

	err = monitorNotificationChannelSlackToResourceData(&nc, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelSlackUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	nc, err := monitorNotificationChannelSlackFromResourceData(d)
	if err != nil {
		return err
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(nc)

	return err
}

func resourceSysdigMonitorNotificationChannelSlackDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteNotificationChannel(id)
}

// Channel type for Notification Channels

func monitorNotificationChannelSlackFromResourceData(d *schema.ResourceData) (nc monitor.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = "SLACK"
	nc.Options.Url = d.Get("url").(string)
	nc.Options.Channel = d.Get("channel").(string)
	return
}

func monitorNotificationChannelSlackToResourceData(nc *monitor.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	d.Set("url", nc.Options.Url)
	d.Set("channel", nc.Options.Channel)
	return
}
