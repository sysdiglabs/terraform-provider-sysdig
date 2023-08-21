package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigMonitorNotificationChannelMSTeams() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigMonitorNotificationChannelMSTeamsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceSysdigMonitorNotificationChannelMSTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelByName(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = monitorNotificationChannelMSTeamsToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(nc.ID))

	return nil
}
