package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
)

func resourceSysdigMonitorNotificationChannelPagerduty() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigMonitorNotificationChannelPagerdutyCreate,
		Update: resourceSysdigMonitorNotificationChannelPagerdutyUpdate,
		Read:   resourceSysdigMonitorNotificationChannelPagerdutyRead,
		Delete: resourceSysdigMonitorNotificationChannelPagerdutyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceSysdigMonitorNotificationChannelPagerdutyCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	notificationChannel, err := monitorNotificationChannelPagerdutyFromResourceData(d)
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
func resourceSysdigMonitorNotificationChannelPagerdutyRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(id)

	if err != nil {
		d.SetId("")
	}

	err = monitorNotificationChannelPagerdutyToResourceData(&nc, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelPagerdutyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	nc, err := monitorNotificationChannelPagerdutyFromResourceData(d)
	if err != nil {
		return err
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(nc)

	return err
}

func resourceSysdigMonitorNotificationChannelPagerdutyDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteNotificationChannel(id)
}

// Channel type for Notification Channels

func monitorNotificationChannelPagerdutyFromResourceData(d *schema.ResourceData) (nc monitor.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_PAGERDUTY
	nc.Options.Account = d.Get("account").(string)
	nc.Options.ServiceKey = d.Get("service_key").(string)
	nc.Options.ServiceName = d.Get("service_name").(string)

	return
}

func monitorNotificationChannelPagerdutyToResourceData(nc *monitor.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	d.Set("account", nc.Options.Account)
	d.Set("service_key", nc.Options.ServiceKey)
	d.Set("service_name", nc.Options.ServiceName)
	return
}
