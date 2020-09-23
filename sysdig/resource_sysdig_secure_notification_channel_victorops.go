package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureNotificationChannelVictorOps() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigSecureNotificationChannelVictorOpsCreate,
		Update: resourceSysdigSecureNotificationChannelVictorOpsUpdate,
		Read:   resourceSysdigSecureNotificationChannelVictorOpsRead,
		Delete: resourceSysdigSecureNotificationChannelVictorOpsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createSecureNotificationChannelSchema(map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routing_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceSysdigSecureNotificationChannelVictorOpsCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	notificationChannel, err := secureNotificationChannelVictorOpsFromResourceData(d)
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
func resourceSysdigSecureNotificationChannelVictorOpsRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(id)

	if err != nil {
		d.SetId("")
	}

	err = secureNotificationChannelVictorOpsToResourceData(&nc, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceSysdigSecureNotificationChannelVictorOpsUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	nc, err := secureNotificationChannelVictorOpsFromResourceData(d)
	if err != nil {
		return err
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(nc)

	return err
}

func resourceSysdigSecureNotificationChannelVictorOpsDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteNotificationChannel(id)
}

func secureNotificationChannelVictorOpsFromResourceData(d *schema.ResourceData) (nc secure.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_VICTOROPS
	nc.Options.APIKey = d.Get("api_key").(string)
	nc.Options.RoutingKey = d.Get("routing_key").(string)
	return
}

func secureNotificationChannelVictorOpsToResourceData(nc *secure.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	d.Set("api_key", nc.Options.APIKey)
	d.Set("routing_key", nc.Options.RoutingKey)
	return
}
