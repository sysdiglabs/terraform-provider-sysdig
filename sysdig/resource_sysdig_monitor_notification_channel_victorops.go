package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigMonitorNotificationChannelVictorOps() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorNotificationChannelVictorOpsCreate,
		UpdateContext: resourceSysdigMonitorNotificationChannelVictorOpsUpdate,
		ReadContext:   resourceSysdigMonitorNotificationChannelVictorOpsRead,
		DeleteContext: resourceSysdigMonitorNotificationChannelVictorOpsDelete,
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

func resourceSysdigMonitorNotificationChannelVictorOpsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := monitorNotificationChannelVictorOpsFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigMonitorNotificationChannelVictorOpsRead(ctx, d, meta)
}

func resourceSysdigMonitorNotificationChannelVictorOpsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
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

	err = monitorNotificationChannelVictorOpsToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelVictorOpsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := monitorNotificationChannelVictorOpsFromResourceData(d, teamID)
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

func resourceSysdigMonitorNotificationChannelVictorOpsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func monitorNotificationChannelVictorOpsFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return
	}

	nc.Type = notificationChannelTypeVictorOps
	nc.Options.APIKey = d.Get("api_key").(string)
	nc.Options.RoutingKey = d.Get("routing_key").(string)
	return
}

func monitorNotificationChannelVictorOpsToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("api_key", nc.Options.APIKey)
	_ = d.Set("routing_key", nc.Options.RoutingKey)

	return
}
