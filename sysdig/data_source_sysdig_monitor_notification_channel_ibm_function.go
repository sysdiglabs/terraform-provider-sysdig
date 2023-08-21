package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigMonitorNotificationChannelIBMFunction() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigMonitorNotificationChannelIBMFunctionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"ibm_function_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_data": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"iam_api_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"whisk_auth_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceSysdigMonitorNotificationChannelIBMFunctionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelByName(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = monitorNotificationChannelIBMFunctionToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(nc.ID))

	return nil
}
