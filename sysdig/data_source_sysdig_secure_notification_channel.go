package sysdig

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	NOTIFICATION_CHANNEL_TYPE_EMAIL                    = "EMAIL"
	NOTIFICATION_CHANNEL_TYPE_AMAZON_SNS               = "SNS"
	NOTIFICATION_CHANNEL_TYPE_OPSGENIE                 = "OPSGENIE"
	NOTIFICATION_CHANNEL_TYPE_VICTOROPS                = "VICTOROPS"
	NOTIFICATION_CHANNEL_TYPE_WEBHOOK                  = "WEBHOOK"
	NOTIFICATION_CHANNEL_TYPE_SLACK                    = "SLACK"
	NOTIFICATION_CHANNEL_TYPE_PAGERDUTY                = "PAGER_DUTY"
	NOTIFICATION_CHANNEL_TYPE_MS_TEAMS                 = "MS_TEAMS"
	NOTIFICATION_CHANNEL_TYPE_GCHAT                    = "GCHAT"
	NOTIFICATION_CHANNEL_TYPE_PROMETHEUS_ALERT_MANAGER = "PROMETHEUS_ALERT_MANAGER"
	NOTIFICATION_CHANNEL_TYPE_TEAM_EMAIL               = "TEAM_EMAIL"
	NOTIFICATION_CHANNEL_TYPE_CUSTOM_WEBHOOK           = "POWER_WEBHOOK"
	NOTIFICATION_CHANNEL_TYPE_IBM_EVENT_NOTIFICATION   = "IBM_EVENT_NOTIFICATIONS"
	NOTIFICATION_CHANNEL_TYPE_IBM_FUNCTION             = "IBM_FUNCTION"

	NOTIFICATION_CHANNEL_TYPE_SLACK_TEMPLATE_KEY_V1    = "SLACK_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v1"
	NOTIFICATION_CHANNEL_TYPE_SLACK_TEMPLATE_KEY_V2    = "SLACK_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v2"
	NOTIFICATION_CHANNEL_TYPE_MS_TEAMS_TEMPLATE_KEY_V1 = "MS_TEAMS_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v1"
	NOTIFICATION_CHANNEL_TYPE_MS_TEAMS_TEMPLATE_KEY_V2 = "MS_TEAMS_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v2"

	NOTIFICATION_CHANNEL_SECURE_EVENT_NOTIFICATION_CONTENT_SECTION = "SECURE_EVENT_NOTIFICATION_CONTENT"
)

func dataSourceSysdigSecureNotificationChannel() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigNotificationChannelRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		DeprecationMessage: "The sysdig_secure_notification_channel data source will be replaced in the next version, " +
			"and split out into different data sources, depending on the type of the notification channel.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recipients": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topics": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routing_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"channel": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notify_when_ok": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_when_resolved": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"send_test_notification": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigNotificationChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelByName(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(nc.ID))
	_ = d.Set("version", nc.Version)
	_ = d.Set("name", nc.Name)
	_ = d.Set("enabled", nc.Enabled)
	_ = d.Set("type", nc.Type)
	_ = d.Set("recipients", strings.Join(nc.Options.EmailRecipients, ","))
	_ = d.Set("topics", strings.Join(nc.Options.SnsTopicARNs, ","))
	_ = d.Set("api_key", nc.Options.APIKey)
	_ = d.Set("url", nc.Options.Url)
	_ = d.Set("channel", nc.Options.Channel)
	_ = d.Set("account", nc.Options.Account)
	_ = d.Set("service_key", nc.Options.ServiceKey)
	_ = d.Set("service_name", nc.Options.ServiceName)
	_ = d.Set("routing_key", nc.Options.RoutingKey)
	_ = d.Set("notify_when_ok", nc.Options.NotifyOnOk)
	_ = d.Set("notify_when_resolved", nc.Options.NotifyOnResolve)
	_ = d.Set("send_test_notification", nc.Options.SendTestNotification)

	// When we receive a notification channel of type OpsGenie,
	// the API sends us the URL, but we are configuring the API
	// key in the file, so terraform identifies this as a change in
	// the resource and tries to update it remotely even if it
	// didn't change at all.
	// We need to extract the key from the url the API gives us
	// to avoid this Terraform's behaviour.
	if nc.Type == NOTIFICATION_CHANNEL_TYPE_OPSGENIE {
		regex, err := regexp.Compile("apiKey=(.*)?$")
		if err != nil {
			return diag.FromErr(err)
		}
		key := regex.FindStringSubmatch(nc.Options.Url)[1]
		_ = d.Set("api_key", key)
		_ = d.Set("url", "")

	}
	return nil
}
