package sysdig

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	NOTIFICATION_CHANNEL_TYPE_EMAIL      = "EMAIL"
	NOTIFICATION_CHANNEL_TYPE_AMAZON_SNS = "SNS"
	NOTIFICATION_CHANNEL_TYPE_OPSGENIE   = "OPSGENIE"
	NOTIFICATION_CHANNEL_TYPE_VICTOROPS  = "VICTOROPS"
	NOTIFICATION_CHANNEL_TYPE_WEBHOOK    = "WEBHOOK"
	NOTIFICATION_CHANNEL_TYPE_SLACK      = "SLACK"
	NOTIFICATION_CHANNEL_TYPE_PAGERDUTY  = "PAGER_DUTY"
)

func dataSourceSysdigSecureNotificationChannel() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Read: dataSourceSysdigNotificationChannelRead,

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
func dataSourceSysdigNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	nc, err := client.GetNotificationChannelByName(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(nc.ID))
	d.Set("version", nc.Version)
	d.Set("name", nc.Name)
	d.Set("enabled", nc.Enabled)
	d.Set("type", nc.Type)
	d.Set("recipients", strings.Join(nc.Options.EmailRecipients, ","))
	d.Set("topics", strings.Join(nc.Options.SnsTopicARNs, ","))
	d.Set("api_key", nc.Options.APIKey)
	d.Set("url", nc.Options.Url)
	d.Set("channel", nc.Options.Channel)
	d.Set("account", nc.Options.Account)
	d.Set("service_key", nc.Options.ServiceKey)
	d.Set("service_name", nc.Options.ServiceName)
	d.Set("routing_key", nc.Options.RoutingKey)
	d.Set("notify_when_ok", nc.Options.NotifyOnOk)
	d.Set("notify_when_resolved", nc.Options.NotifyOnResolve)
	d.Set("send_test_notification", nc.Options.SendTestNotification)

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
			return err
		}
		key := regex.FindStringSubmatch(nc.Options.Url)[1]
		d.Set("api_key", key)
		d.Set("url", "")
	}
	return nil
}
