package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureNotificationChannelEmail() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigSecureNotificationChannelEmailCreate,
		Update: resourceSysdigSecureNotificationChannelEmailUpdate,
		Read:   resourceSysdigSecureNotificationChannelEmailRead,
		Delete: resourceSysdigSecureNotificationChannelEmailDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createSecureNotificationChannelSchema(map[string]*schema.Schema{
			"recipients": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		}),
	}
}

func resourceSysdigSecureNotificationChannelEmailCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	notificationChannel, err := secureNotificationChannelEmailFromResourceData(d)
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
func resourceSysdigSecureNotificationChannelEmailRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(id)

	if err != nil {
		d.SetId("")
	}

	err = secureNotificationChannelEmailToResourceData(&nc, d)
	if err != nil {
		return err
	}

	return nil
}

func resourceSysdigSecureNotificationChannelEmailUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	nc, err := secureNotificationChannelEmailFromResourceData(d)
	if err != nil {
		return err
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(nc)

	return err
}

func resourceSysdigSecureNotificationChannelEmailDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteNotificationChannel(id)
}

// Channel type for Notification Channels

func secureNotificationChannelEmailFromResourceData(d *schema.ResourceData) (nc secure.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_EMAIL
	nc.Options.EmailRecipients = cast.ToStringSlice(d.Get("recipients").(*schema.Set).List())
	return
}

func secureNotificationChannelEmailToResourceData(nc *secure.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	d.Set("recipients", nc.Options.EmailRecipients)
	return
}
