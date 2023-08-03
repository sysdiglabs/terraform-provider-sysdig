//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_common || tf_acc_ibm_monitor || tf_acc_ibm_common

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelCustomWebhook(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: sysdigOrIBMMonitorPreCheck(t),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorNotificationChannelCustomWebhookWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_custom_webhook.sample-custom-webhook1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelCustomWebhookWithNameWithAdditionalheaders(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_custom_webhook.sample-custom-webhook2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelCustomWebhookSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_custom_webhook.sample-custom-webhook3",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelCustomWebhookSharedWithAllowInsecureConnections(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_custom_webhook.sample-custom-webhook4",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelCustomWebhookSharedWithAdditionalHeaders(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_custom_webhook.sample-custom-webhook5",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelCustomWebhookWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_custom_webhook" "sample-custom-webhook1" {
	name = "Example Channel %s - Custom Webhook"
	enabled = true
	url = "https://example.com/"
	http_method = "POST"
	template = "{\n  \"code\": \"incident\",\n  \"alert\": \"{{@alert_name}}\"\n}"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func monitorNotificationChannelCustomWebhookWithNameWithAdditionalheaders(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_custom_webhook" "sample-custom-webhook2" {
		name = "Example Channel %s - Custom Webhook With Additional Headers"
		enabled = true
		url = "https://example.com/"
		http_method = "POST"
		template = "{\n  \"code\": \"incident\",\n  \"alert\": \"{{@alert_name}}\"\n}"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		additional_headers = {
			"Webhook-Header": "TestHeader"
		}
	}`, name)
}

func monitorNotificationChannelCustomWebhookSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_custom_webhook" "sample-custom-webhook3" {
		name = "Example Channel %s - Custom Webhook With Additional Headers"
		share_with_current_team = true
		enabled = true
		url = "https://example.com/"
		http_method = "POST"
		template = "{\n  \"code\": \"incident\",\n  \"alert\": \"{{@alert_name}}\"\n}"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
	}`, name)
}

func monitorNotificationChannelCustomWebhookSharedWithAllowInsecureConnections(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_custom_webhook" "sample-custom-webhook4" {
		name = "Example Channel %s - Custom Webhook With Additional Headers"
		enabled = true
		url = "https://example.com/"
		http_method = "POST"
		template = "{\n  \"code\": \"incident\",\n  \"alert\": \"{{@alert_name}}\"\n}"
		allow_insecure_connections = true
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
	}`, name)
}

func monitorNotificationChannelCustomWebhookSharedWithAdditionalHeaders(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_custom_webhook" "sample-custom-webhook5" {
		name = "Example Channel %s - Custom Webhook With Additional Headers"
		enabled = true
		url = "https://example.com/"
		http_method = "POST"
		template = "{\n  \"code\": \"incident\",\n  \"alert\": \"{{@alert_name}}\"\n}"
		additional_headers = {
			"Webhook-Header": "TestHeader"
		}
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
	}`, name)
}
