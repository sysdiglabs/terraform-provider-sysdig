package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigMonitorNotificationChannelSlack() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorNotificationChannelSlackCreate,
		UpdateContext: resourceSysdigMonitorNotificationChannelSlackUpdate,
		ReadContext:   resourceSysdigMonitorNotificationChannelSlackRead,
		DeleteContext: resourceSysdigMonitorNotificationChannelSlackDelete,
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
			"url": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"channel": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private_channel": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"private_channel_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"show_section_runbook_links": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_section_event_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_section_user_defined_content": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_section_notification_chart": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_section_dashboard_links": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_section_alert_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_section_capturing_information": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		}),
	}
}

func resourceSysdigMonitorNotificationChannelSlackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := monitorNotificationChannelSlackFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigMonitorNotificationChannelSlackRead(ctx, d, meta)
}

func resourceSysdigMonitorNotificationChannelSlackRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelById(ctx, id)
	if err != nil {
		if err == v2.NotificationChannelNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = monitorNotificationChannelSlackToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelSlackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := monitorNotificationChannelSlackFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	nc.Version = d.Get("version").(int)
	nc.ID, err = strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateNotificationChannel(ctx, nc)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigMonitorNotificationChannelSlackRead(ctx, d, meta)

	return nil
}

func resourceSysdigMonitorNotificationChannelSlackDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteNotificationChannel(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func monitorNotificationChannelSlackFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = monitorNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_SLACK
	nc.Options.Url = d.Get("url").(string)
	nc.Options.Channel = d.Get("channel").(string)
	nc.Options.PrivateChannel = d.Get("is_private_channel").(bool)
	nc.Options.PrivateChannelUrl = d.Get("private_channel_url").(string)
	nc.Options.TemplateConfiguration = []v2.NotificationChannelTemplateConfiguration{
		{
			TemplateKey: "SLACK_MONITOR_ALERT_NOTIFICATION_TEMPLATE_METADATA_v1",
			TemplateConfigurationSections: []v2.NotificationChannelTemplateConfigurationSection{
				{
					SectionName: "MONITOR_ALERT_NOTIFICATION_HEADER",
					ShouldShow:  true,
				},
				{
					SectionName: "MONITOR_ALERT_NOTIFICATION_RUNBOOK_LINKS",
					ShouldShow:  d.Get("show_section_runbook_links").(bool),
				},
				{
					SectionName: "MONITOR_ALERT_NOTIFICATION_EVENT_DETAILS",
					ShouldShow:  d.Get("show_section_event_details").(bool),
				},
				{
					SectionName: "MONITOR_ALERT_NOTIFICATION_USER_DEFINED_CONTENT",
					ShouldShow:  d.Get("show_section_user_defined_content").(bool),
				},
				{
					SectionName: "MONITOR_ALERT_NOTIFICATION_CHART",
					ShouldShow:  d.Get("show_section_notification_chart").(bool),
				},
				{
					SectionName: "MONITOR_ALERT_NOTIFICATION_DASHBOARD_LINKS",
					ShouldShow:  d.Get("show_section_dashboard_links").(bool),
				},
				{
					SectionName: "MONITOR_ALERT_NOTIFICATION_ALERT_DETAILS",
					ShouldShow:  d.Get("show_section_alert_details").(bool),
				},
				{
					SectionName: "MONITOR_ALERT_NOTIFICATION_CAPTURING_INFORMATION",
					ShouldShow:  d.Get("show_section_capturing_information").(bool),
				},
			},
		},
	}

	return
}

func monitorNotificationChannelSlackToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = monitorNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("url", nc.Options.Url)
	_ = d.Set("channel", nc.Options.Channel)
	_ = d.Set("is_private_channel", nc.Options.PrivateChannel)
	_ = d.Set("private_channel_url", nc.Options.PrivateChannelUrl)

	runbookLinks := true
	eventDetails := true
	userDefinedContent := true
	notificationChart := true
	dashboardLinks := true
	alertDetails := true
	capturingInformation := true

	if len(nc.Options.TemplateConfiguration) == 1 {
		for _, c := range nc.Options.TemplateConfiguration[0].TemplateConfigurationSections {
			switch c.SectionName {
			case "MONITOR_ALERT_NOTIFICATION_RUNBOOK_LINKS":
				runbookLinks = c.ShouldShow
			case "MONITOR_ALERT_NOTIFICATION_EVENT_DETAILS":
				eventDetails = c.ShouldShow
			case "MONITOR_ALERT_NOTIFICATION_USER_DEFINED_CONTENT":
				userDefinedContent = c.ShouldShow
			case "MONITOR_ALERT_NOTIFICATION_CHART":
				notificationChart = c.ShouldShow
			case "MONITOR_ALERT_NOTIFICATION_DASHBOARD_LINKS":
				dashboardLinks = c.ShouldShow
			case "MONITOR_ALERT_NOTIFICATION_ALERT_DETAILS":
				alertDetails = c.ShouldShow
			case "MONITOR_ALERT_NOTIFICATION_CAPTURING_INFORMATION":
				capturingInformation = c.ShouldShow
			}
		}
	}

	_ = d.Set("show_section_runbook_links", runbookLinks)
	_ = d.Set("show_section_event_details", eventDetails)
	_ = d.Set("show_section_user_defined_content", userDefinedContent)
	_ = d.Set("show_section_notification_chart", notificationChart)
	_ = d.Set("show_section_dashboard_links", dashboardLinks)
	_ = d.Set("show_section_alert_details", alertDetails)
	_ = d.Set("show_section_capturing_information", capturingInformation)

	return
}
