package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureNotificationChannelOpsGenie() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureNotificationChannelOpsGenieCreate,
		UpdateContext: resourceSysdigSecureNotificationChannelOpsGenieUpdate,
		ReadContext:   resourceSysdigSecureNotificationChannelOpsGenieRead,
		DeleteContext: resourceSysdigSecureNotificationChannelOpsGenieDelete,
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
			"api_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "US",
				ValidateFunc: validation.StringInSlice([]string{"US", "EU"}, false),
			},
		}),
	}
}

func resourceSysdigSecureNotificationChannelOpsGenieCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := secureNotificationChannelOpsGenieFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigSecureNotificationChannelOpsGenieRead(ctx, d, meta)
}

func resourceSysdigSecureNotificationChannelOpsGenieRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	err = secureNotificationChannelOpsGenieToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureNotificationChannelOpsGenieUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := secureNotificationChannelOpsGenieFromResourceData(d, teamID)
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

func resourceSysdigSecureNotificationChannelOpsGenieDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func secureNotificationChannelOpsGenieFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_OPSGENIE
	apiKey := d.Get("api_key").(string)
	nc.Options.APIKey = apiKey
	nc.Options.Region = d.Get("region").(string)
	return
}

func secureNotificationChannelOpsGenieToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("api_key", nc.Options.APIKey)
	_ = d.Set("region", nc.Options.Region)

	return
}
