package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureNotificationChannelPagerduty() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureNotificationChannelPagerdutyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createSecureNotificationChannelSchema(map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceSysdigSecureNotificationChannelPagerdutyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))

	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelByName(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = secureNotificationChannelPagerdutyToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(nc.ID))

	return nil
}
