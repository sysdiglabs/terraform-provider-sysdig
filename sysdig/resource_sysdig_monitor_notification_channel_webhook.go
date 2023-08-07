package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigMonitorNotificationChannelWebhook() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorNotificationChannelWebhookCreate,
		UpdateContext: resourceSysdigMonitorNotificationChannelWebhookUpdate,
		ReadContext:   resourceSysdigMonitorNotificationChannelWebhookRead,
		DeleteContext: resourceSysdigMonitorNotificationChannelWebhookDelete,
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
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"additional_headers": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"allow_insecure_connections": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		}),
	}
}

func resourceSysdigMonitorNotificationChannelWebhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorNotificationChannelClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := monitorNotificationChannelWebhookFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigMonitorNotificationChannelWebhookRead(ctx, d, meta)
}

func resourceSysdigMonitorNotificationChannelWebhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(ctx, id)
	if err != nil {
		if err == v2.NotificationChannelNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = monitorNotificationChannelWebhookToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := monitorNotificationChannelWebhookFromResourceData(d, teamID)
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

func resourceSysdigMonitorNotificationChannelWebhookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
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

func monitorNotificationChannelWebhookFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_WEBHOOK
	nc.Options.Url = d.Get("url").(string)
	nc.Options.AdditionalHeaders = d.Get("additional_headers").(map[string]interface{})
	allowInsecureConnections := d.Get("allow_insecure_connections").(bool)
	nc.Options.AllowInsecureConnections = &allowInsecureConnections
	return
}

func monitorNotificationChannelWebhookToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("url", nc.Options.Url)
	_ = d.Set("additional_headers", nc.Options.AdditionalHeaders)
	if nc.Options.AllowInsecureConnections != nil {
		_ = d.Set("allow_insecure_connections", *nc.Options.AllowInsecureConnections)
	}

	return
}
