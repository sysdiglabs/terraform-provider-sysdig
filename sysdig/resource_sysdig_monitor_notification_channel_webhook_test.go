//go:build tf_acc_sysdig || tf_acc_ibm

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelWebhook(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			ibmMonitor := os.Getenv("SYSDIG_IBM_MONITOR_API_KEY")
			if monitor != "" || ibmMonitor != "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN or SYSDIG_IBM_MONITOR_API_KEY must be set for acceptance tests")
			}
		},
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
				ResourceName:      "sysdig_monitor_notification_channel_webhook.sample-webhook",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelWebhookWithNameWithAdditionalheaders(rText()),
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
