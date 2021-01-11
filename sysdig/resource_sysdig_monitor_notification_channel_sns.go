package sysdig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
)

func resourceSysdigMonitorNotificationChannelSNS() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorNotificationChannelSNSCreate,
		UpdateContext: resourceSysdigMonitorNotificationChannelSNSUpdate,
		ReadContext:   resourceSysdigMonitorNotificationChannelSNSRead,
		DeleteContext: resourceSysdigMonitorNotificationChannelSNSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

func resourceSysdigMonitorNotificationChannelSNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := monitorNotificationChannelSNSFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))
	d.Set("version", notificationChannel.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigMonitorNotificationChannelSNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(ctx, id)

	if err != nil {
		d.SetId("")
	}

	err = monitorNotificationChannelSNSToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelSNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := monitorNotificationChannelSNSFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(ctx, nc)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelSNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteNotificationChannel(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
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
