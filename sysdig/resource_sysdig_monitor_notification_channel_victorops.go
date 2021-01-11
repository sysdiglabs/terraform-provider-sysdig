package sysdig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
)

func resourceSysdigMonitorNotificationChannelVictorOps() *schema.Resource {
	timeout := 30 * time.Second

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

func resourceSysdigMonitorNotificationChannelVictorOpsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := monitorNotificationChannelVictorOpsFromResourceData(d)
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
func resourceSysdigMonitorNotificationChannelVictorOpsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(ctx, id)

	if err != nil {
		d.SetId("")
	}

	err = monitorNotificationChannelVictorOpsToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelVictorOpsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := monitorNotificationChannelVictorOpsFromResourceData(d)
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

func resourceSysdigMonitorNotificationChannelVictorOpsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func monitorNotificationChannelVictorOpsFromResourceData(d *schema.ResourceData) (nc monitor.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_VICTOROPS
	nc.Options.APIKey = d.Get("api_key").(string)
	nc.Options.RoutingKey = d.Get("routing_key").(string)
	return
}

func monitorNotificationChannelVictorOpsToResourceData(nc *monitor.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	d.Set("api_key", nc.Options.APIKey)
	d.Set("routing_key", nc.Options.RoutingKey)
	return
}
