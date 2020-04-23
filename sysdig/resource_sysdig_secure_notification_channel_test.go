package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccNotificationChannel(t *testing.T) {
	//var ncBefore, ncAfter secure.NotificationChannel

	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: notificationChannelEmailWithName(rText()),
			},
			{
				Config: notificationChannelAmazonSNSWithName(rText()),
			},
			{
				Config: notificationChannelOpsGenieWithName(rText()),
			},
			{
				Config: notificationChannelVictorOpsWithName(rText()),
			},
			{
				Config: notificationChannelWebhookWithName(rText()),
			},
			{
				Config: notificationChannelSlackWithName(rText()),
			},
			{
				Config: notificationChannelPagerdutyWithName(rText()),
			},
		},
	})
}

func notificationChannelEmailWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample_email" {
	name = "%s"
	enabled = true
	type = "EMAIL"
	recipients = "root@localhost.com"
	notify_when_ok = false
	notify_when_resolved = false
}`, name)
}

func notificationChannelAmazonSNSWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample-amazon-sns" {
	name = "Example Channel %s - Amazon SNS"
	enabled = true
	type = "SNS"
	topics = "arn:aws:sns:us-east-1:273489009834:my-alerts,arn:aws:sns:us-east-1:279948934544:my-alerts2"
	notify_when_ok = false
	notify_when_resolved = false
}`, name)
}

func notificationChannelVictorOpsWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample-victorops" {
	name = "Example Channel %s - VictorOps"
	enabled = true
	type = "VICTOROPS"
	api_key = "1234342-4234243-4234-2"
	routing_key = "My team"
	notify_when_ok = false
	notify_when_resolved = false
}`, name)
}

func notificationChannelOpsGenieWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample-opsgenie" {
	name = "Example Channel %s - OpsGenie"
	enabled = true
	type = "OPSGENIE"
	api_key = "2349324-342354353-5324-23"
	notify_when_ok = false
	notify_when_resolved = false
}`, name)
}

func notificationChannelWebhookWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample-webhook" {
	name = "Example Channel %s - Webhook"
	enabled = true
	type = "WEBHOOK"
	url = "localhost:8080"
	notify_when_ok = false
	notify_when_resolved = false
}`, name)
}

func notificationChannelSlackWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample-slack" {
	name = "Example Channel %s - Slack"
	enabled = true
	type = "SLACK"
	url = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}

func notificationChannelPagerdutyWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample-pagerduty" {
	name = "Example Channel %s - Pagerduty"
	enabled = true
	type = "PAGER_DUTY"
	account = "account"
	service_key = "XXXXXXXXXX"
	service_name = "sysdig"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}
