package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureNotificationChannelWebhook() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureNotificationChannelWebhookCreate,
		UpdateContext: resourceSysdigSecureNotificationChannelWebhookUpdate,
		ReadContext:   resourceSysdigSecureNotificationChannelWebhookRead,
		DeleteContext: resourceSysdigSecureNotificationChannelWebhookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createSecureNotificationChannelSchema(map[string]*schema.Schema{
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
			"custom_data": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		}),
	}
}

func resourceSysdigSecureNotificationChannelWebhookCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := secureNotificationChannelWebhookFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigSecureNotificationChannelWebhookRead(ctx, d, meta)
}

func resourceSysdigSecureNotificationChannelWebhookRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelByID(ctx, id)
	if err != nil {
		if err == v2.ErrNotificationChannelNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = secureNotificationChannelWebhookToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureNotificationChannelWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := secureNotificationChannelWebhookFromResourceData(d, teamID)
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

func resourceSysdigSecureNotificationChannelWebhookDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
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

func secureNotificationChannelWebhookFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return nc, err
	}

	nc.Type = notificationChannelTypeWebhook
	nc.Options.URL = d.Get("url").(string)
	nc.Options.AdditionalHeaders = d.Get("additional_headers").(map[string]any)
	nc.Options.CustomData = d.Get("custom_data").(map[string]any)
	allowInsecureConnections := d.Get("allow_insecure_connections").(bool)
	nc.Options.AllowInsecureConnections = &allowInsecureConnections
	return nc, err
}

func secureNotificationChannelWebhookToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return err
	}

	_ = d.Set("url", nc.Options.URL)
	_ = d.Set("additional_headers", nc.Options.AdditionalHeaders)
	_ = d.Set("custom_data", nc.Options.CustomData)
	if nc.Options.AllowInsecureConnections != nil {
		_ = d.Set("allow_insecure_connections", *nc.Options.AllowInsecureConnections)
	}

	return err
}
