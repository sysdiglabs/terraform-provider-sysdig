package sysdig

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureNotificationChannel() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigNotificationChannelCreate,
		Update: resourceSysdigNotificationChannelUpdate,
		Read:   resourceSysdigNotificationChannelRead,
		Delete: resourceSysdigNotificationChannelDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"recipients": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topics": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routing_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_when_ok": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"notify_when_resolved": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSysdigNotificationChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	notificationChannel, err := notificationChannelFromResourceData(d)
	if err != nil {
		return err
	}

	notificationChannel, err = client.CreateNotificationChannel(notificationChannel)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))
	d.Set("version", notificationChannel.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(id)

	if err != nil {
		d.SetId("")
	}

	d.Set("version", nc.Version)
	d.Set("name", nc.Name)
	d.Set("enabled", nc.Enabled)
	d.Set("type", nc.Type)
	d.Set("recipients", nc.Options.EmailRecipients)
	d.Set("topics", nc.Options.SnsTopicARNs)
	d.Set("api_key", nc.Options.APIKey)
	d.Set("url", nc.Options.Url)
	d.Set("channel", nc.Options.Channel)
	d.Set("account", nc.Options.Account)
	d.Set("service_key", nc.Options.ServiceKey)
	d.Set("service_name", nc.Options.ServiceName)
	d.Set("routing_key", nc.Options.RoutingKey)
	d.Set("notify_when_ok", nc.Options.NotifyOnOk)
	d.Set("notify_when_resolved", nc.Options.NotifyOnResolve)

	// When we receive a notification channel of type OpsGenie,
	// the API sends us the URL, but we are configuring the API
	// key in the file, so terraform identifies this as a change in
	// the resource and tries to update it remotely even if it
	// didn't change at all.
	// We need to extract the key from the url the API gives us
	// to avoid this Terraform's behaviour.
	if nc.Type == opsgenie {
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

func resourceSysdigNotificationChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	nc, err := notificationChannelFromResourceData(d)
	if err != nil {
		return err
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(nc)

	return err
}

func resourceSysdigNotificationChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteNotificationChannel(id)
}

// Channel type for Notification Channels
const (
	email     = "EMAIL"
	amazonSns = "SNS"
	opsgenie  = "OPSGENIE"
	victorops = "VICTOROPS"
	webhook   = "WEBHOOK"
	slack     = "SLACK"
	pagerduty = "PAGER_DUTY"
)

func notificationChannelFromResourceData(d *schema.ResourceData) (nc secure.NotificationChannel, err error) {

	channelType := strings.ToUpper(d.Get("type").(string))

	nc = secure.NotificationChannel{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
		Type:    channelType,
		Options: secure.NotificationChannelOptions{
			NotifyOnOk:      d.Get("notify_when_ok").(bool),
			NotifyOnResolve: d.Get("notify_when_resolved").(bool),
		},
	}

	fieldNotSetError := "the '%s' field must be set when the type of the notification channel is %s"

	// Retrieve the special options for each type
	switch channelType {
	case email:
		if recipients, ok := d.Get("recipients").(string); ok && recipients != "" {
			emails := strings.Split(recipients, ",")
			for _, email := range emails {
				// We need to trim the emails or the API will not accept them
				nc.Options.EmailRecipients = append(nc.Options.EmailRecipients, strings.TrimSpace(email))
			}
		} else {
			err = fmt.Errorf(fieldNotSetError, "recipients", channelType)
			return
		}
	case amazonSns:
		if snsTopics, ok := d.Get("topics").(string); ok && snsTopics != "" {
			topics := strings.Split(snsTopics, ",")
			for _, topic := range topics {
				// We need to trim the topics or the API will not accept them
				nc.Options.SnsTopicARNs = append(nc.Options.SnsTopicARNs, strings.TrimSpace(topic))
			}
		} else {
			err = fmt.Errorf(fieldNotSetError, "topics", channelType)
			return
		}
	case victorops:
		if apiKey, ok := d.Get("api_key").(string); ok && apiKey != "" {
			nc.Options.APIKey = apiKey
		} else {
			err = fmt.Errorf(fieldNotSetError, "api_key", channelType)
			return
		}

		if routingKey, ok := d.Get("routing_key").(string); ok && routingKey != "" {
			nc.Options.RoutingKey = routingKey
		} else {
			err = fmt.Errorf(fieldNotSetError, "routing_key", channelType)
			return
		}
	case opsgenie:
		if apiKey, ok := d.Get("api_key").(string); ok && apiKey != "" {
			nc.Options.Url = fmt.Sprintf("https://api.opsgenie.com/v1/json/sysdigcloud?apiKey=%s", apiKey)
		} else {
			err = fmt.Errorf(fieldNotSetError, "api_key", channelType)
			return
		}
	case webhook:
		if url, ok := d.Get("url").(string); ok && url != "" {
			nc.Options.Url = url
		} else {
			err = fmt.Errorf(fieldNotSetError, "url", channelType)
			return
		}
	case slack:
		if url, ok := d.Get("url").(string); ok && url != "" {
			nc.Options.Url = url
		} else {
			err = fmt.Errorf(fieldNotSetError, "url", channelType)
			return
		}
		if channel, ok := d.Get("channel").(string); ok && channel != "" {
			nc.Options.Channel = channel
		} else {
			err = fmt.Errorf(fieldNotSetError, "channel", channelType)
			return
		}
	case pagerduty:
		if account, ok := d.Get("account").(string); ok && account != "" {
			nc.Options.Account = account
		} else {
			err = fmt.Errorf(fieldNotSetError, "account", channelType)
			return
		}
		if serviceKey, ok := d.Get("service_key").(string); ok && serviceKey != "" {
			nc.Options.ServiceKey = serviceKey
		} else {
			err = fmt.Errorf(fieldNotSetError, "service_key", channelType)
			return
		}
		if serviceName, ok := d.Get("service_name").(string); ok && serviceName != "" {
			nc.Options.ServiceName = serviceName
		} else {
			err = fmt.Errorf(fieldNotSetError, "service_name", channelType)
			return
		}
	default:
		validChannelTypes := []string{email, amazonSns, opsgenie, victorops, webhook, slack, pagerduty}
		err = fmt.Errorf("error type not recognized, must be one of the following: %s", validChannelTypes)
		return
	}

	return
}
