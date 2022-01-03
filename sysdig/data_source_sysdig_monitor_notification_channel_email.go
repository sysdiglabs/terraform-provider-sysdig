package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigMonitorNotificationChannelEmail() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigMonitorNotificationChannelEmailRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"recipients": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}

func dataSourceSysdigMonitorNotificationChannelEmailRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigMonitorClient()

	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelByName(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = monitorNotificationChannelEmailToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(nc.ID))

	return nil
}
