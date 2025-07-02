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
	notificationChannelTypeEmail                  = "EMAIL"
	notificationChannelTypeAmazonSNS              = "SNS"
	notificationChannelTypeOpsGenie               = "OPSGENIE"
	notificationChannelTypeVictorOps              = "VICTOROPS"
	notificationChannelTypeWebhook                = "WEBHOOK"
	notificationChannelTypeSlack                  = "SLACK"
	notificationChannelTypePagerduty              = "PAGER_DUTY"
	notificationChannelTypeMSTeams                = "MS_TEAMS"
	notificationChannelTypeGChat                  = "GCHAT"
	notificationChannelTypePrometheusAlertManager = "PROMETHEUS_ALERT_MANAGER"
	notificationChannelTypeTeamEmail              = "TEAM_EMAIL"
	notificationChannelTypeCustomWebhook          = "POWER_WEBHOOK"
	notificationChannelTypeIBMEventNotification   = "IBM_EVENT_NOTIFICATIONS"

	notificationChannelTypeSlackTemplateKeyV1   = "SLACK_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v1"
	notificationChannelTypeSlackTemplateKeyV2   = "SLACK_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v2"
	notificationChannelTypeMSTeamsTemplateKeyV1 = "MS_TEAMS_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v1"
	notificationChannelTypeMSTeamsTemplateKeyV2 = "MS_TEAMS_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v2"

	notificationChannelSecureEventNotificationContentSection = "SECURE_EVENT_NOTIFICATION_CONTENT"
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
func dataSourceSysdigNotificationChannelRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	_ = d.Set("url", nc.Options.URL)
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
	if nc.Type == notificationChannelTypeOpsGenie {
		regex, err := regexp.Compile("apiKey=(.*)?$")
		if err != nil {
			return diag.FromErr(err)
		}
		key := regex.FindStringSubmatch(nc.Options.URL)[1]
		_ = d.Set("api_key", key)
		_ = d.Set("url", "")

	}
	return nil
}
