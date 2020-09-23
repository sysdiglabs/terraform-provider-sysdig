package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
)

func resourceSysdigMonitorNotificationChannelSNS() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigMonitorNotificationChannelSNSCreate,
		Update: resourceSysdigMonitorNotificationChannelSNSUpdate,
		Read:   resourceSysdigMonitorNotificationChannelSNSRead,
		Delete: resourceSysdigMonitorNotificationChannelSNSDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"topics": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		}),
	}
}

func resourceSysdigMonitorNotificationChannelSNSCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	notificationChannel, err := monitorNotificationChannelSNSFromResourceData(d)
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
func resourceSysdigMonitorNotificationChannelSNSRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(id)

	if err != nil {
		d.SetId("")
	}

	err = monitorNotificationChannelSNSToResourceData(&nc, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelSNSUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	nc, err := monitorNotificationChannelSNSFromResourceData(d)
	if err != nil {
		return err
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(nc)

	return err
}

func resourceSysdigMonitorNotificationChannelSNSDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteNotificationChannel(id)
}

// Channel type for Notification Channels

func monitorNotificationChannelSNSFromResourceData(d *schema.ResourceData) (nc monitor.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_AMAZON_SNS
	nc.Options.SnsTopicARNs = cast.ToStringSlice(d.Get("topics").(*schema.Set).List())
	return
}

func monitorNotificationChannelSNSToResourceData(nc *monitor.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	d.Set("topics", nc.Options.SnsTopicARNs)
	return
}
