package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureNotificationChannelPagerduty() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureNotificationChannelPagerdutyCreate,
		UpdateContext: resourceSysdigSecureNotificationChannelPagerdutyUpdate,
		ReadContext:   resourceSysdigSecureNotificationChannelPagerdutyRead,
		DeleteContext: resourceSysdigSecureNotificationChannelPagerdutyDelete,
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

func resourceSysdigSecureNotificationChannelPagerdutyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := secureNotificationChannelPagerdutyFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigSecureNotificationChannelPagerdutyRead(ctx, d, meta)
}

func resourceSysdigSecureNotificationChannelPagerdutyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
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

	err = secureNotificationChannelPagerdutyToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureNotificationChannelPagerdutyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := secureNotificationChannelPagerdutyFromResourceData(d, teamID)
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

func resourceSysdigSecureNotificationChannelPagerdutyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func secureNotificationChannelPagerdutyFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_PAGERDUTY
	nc.Options.Account = d.Get("account").(string)
	nc.Options.ServiceKey = d.Get("service_key").(string)
	nc.Options.ServiceName = d.Get("service_name").(string)

	return
}

func secureNotificationChannelPagerdutyToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("account", nc.Options.Account)
	_ = d.Set("service_key", nc.Options.ServiceKey)
	_ = d.Set("service_name", nc.Options.ServiceName)

	return
}
