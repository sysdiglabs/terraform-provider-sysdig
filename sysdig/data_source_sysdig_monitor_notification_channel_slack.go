package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigMonitorNotificationChannelSlack() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigMonitorNotificationChannelSlackRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"url": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"channel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_private_channel": {
				Type:     schema.TypeBool,
				Required: false,
			},
			"private_channel_url": {
				Type:     schema.TypeString,
				Required: false,
			},
			"show_section_runbook_links": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_section_event_details": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_section_user_defined_content": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_section_notification_chart": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_section_dashboard_links": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_section_alert_details": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_section_capturing_information": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		}),
	}
}

func dataSourceSysdigMonitorNotificationChannelSlackRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelByName(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = monitorNotificationChannelSlackToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(nc.ID))

	return nil
}
