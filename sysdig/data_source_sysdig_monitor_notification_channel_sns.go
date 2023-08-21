package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigMonitorNotificationChannelSNS() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigMonitorNotificationChannelSNSRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"topics": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}

func dataSourceSysdigMonitorNotificationChannelSNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelByName(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = monitorNotificationChannelSNSToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(nc.ID))

	return nil
}
