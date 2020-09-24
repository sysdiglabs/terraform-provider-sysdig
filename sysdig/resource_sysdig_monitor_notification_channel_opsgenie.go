package sysdig

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
)

func resourceSysdigMonitorNotificationChannelOpsGenie() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigMonitorNotificationChannelOpsGenieCreate,
		Update: resourceSysdigMonitorNotificationChannelOpsGenieUpdate,
		Read:   resourceSysdigMonitorNotificationChannelOpsGenieRead,
		Delete: resourceSysdigMonitorNotificationChannelOpsGenieDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceSysdigMonitorNotificationChannelOpsGenieCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	notificationChannel, err := monitorNotificationChannelOpsGenieFromResourceData(d)
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
func resourceSysdigMonitorNotificationChannelOpsGenieRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(id)

	if err != nil {
		d.SetId("")
	}

	err = monitorNotificationChannelOpsGenieToResourceData(&nc, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelOpsGenieUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	nc, err := monitorNotificationChannelOpsGenieFromResourceData(d)
	if err != nil {
		return err
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(nc)

	return err
}

func resourceSysdigMonitorNotificationChannelOpsGenieDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteNotificationChannel(id)
}

// Channel type for Notification Channels

func monitorNotificationChannelOpsGenieFromResourceData(d *schema.ResourceData) (nc monitor.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_OPSGENIE
	apiKey := d.Get("api_key").(string)
	nc.Options.Url = fmt.Sprintf("https://api.opsgenie.com/v1/json/sysdigcloud?apiKey=%s", apiKey)
	return
}

func monitorNotificationChannelOpsGenieToResourceData(nc *monitor.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	regex, err := regexp.Compile("apiKey=(.*)?$")
	if err != nil {
		return err
	}
	key := regex.FindStringSubmatch(nc.Options.Url)[1]
	d.Set("api_key", key)
	return
}
