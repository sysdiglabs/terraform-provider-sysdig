//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_common || tf_acc_ibm_monitor || tf_acc_ibm_common || tf_acc_onprem_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelWebhook(t *testing.T) {
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
				Config: monitorNotificationChannelWebhookWithName(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_webhook.sample-webhook",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: monitorNotificationChannelWebhookWithNameWithAdditionalheaders(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_webhook.sample-webhook2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: monitorNotificationChannelWebhookSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_webhook.sample-webhook3",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: monitorNotificationChannelWebhookSharedWithAllowInsecureConnections(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_webhook.sample-webhook4",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: monitorNotificationChannelWebhookSharedWithCustomData(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_webhook.sample-webhook5",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
		},
	})
}

func monitorNotificationChannelWebhookWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_webhook" "sample-webhook" {
	name = "Example Channel %s - Webhook"
	enabled = true
	url = "https://example.com/"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func monitorNotificationChannelWebhookWithNameWithAdditionalheaders(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_webhook" "sample-webhook2" {
		name = "Example Channel %s - Webhook With Additional Headers"
		enabled = true
		url = "https://example.com/"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		additional_headers = {
			"Webhook-Header": "TestHeader"
		}
	}`, name)
}

func monitorNotificationChannelWebhookSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_webhook" "sample-webhook3" {
		name = "Example Channel %s - Webhook With Additional Headers"
		share_with_current_team = true
		enabled = true
		url = "https://example.com/"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
	}`, name)
}

func monitorNotificationChannelWebhookSharedWithAllowInsecureConnections(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_webhook" "sample-webhook4" {
		name = "Example Channel %s - Webhook With Allow Insecure Connections"
		enabled = true
		allow_insecure_connections = true
		url = "https://example.com/"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
	}`, name)
}

func monitorNotificationChannelWebhookSharedWithCustomData(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_webhook" "sample-webhook5" {
		name = "Example Channel %s - Webhook With Custom Data"
		enabled = true
		url = "https://example.com/"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		custom_data = {
			"data1": "value1"
			"data2": "value2"
		}
	}`, name)
}
