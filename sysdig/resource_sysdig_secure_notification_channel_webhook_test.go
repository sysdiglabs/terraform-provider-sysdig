//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelWebhook(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelWebhookWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_webhook.sample-webhook",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: secureNotificationChannelWebhookSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_webhook.sample-webhook3",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: secureNotificationChannelWebhookSharedWithAllowInsecureConnections(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_webhook.sample-webhook4",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: secureNotificationChannelWebhookSharedWithCustomData(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_webhook.sample-webhook5",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelWebhookWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_webhook" "sample-webhook" {
	name = "Example Channel %s - Webhook"
	enabled = true
	url = "http://1.1.1.1:8080"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func secureNotificationChannelWebhookSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_notification_channel_webhook" "sample-webhook3" {
		name = "Example Channel %s - Webhook With Additional Headers"
		share_with_current_team = true
		enabled = true
		url = "https://example.com/"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
	}`, name)
}

func secureNotificationChannelWebhookSharedWithAllowInsecureConnections(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_notification_channel_webhook" "sample-webhook4" {
		name = "Example Channel %s - Webhook With Allow Insecure Connections"
		enabled = true
		allow_insecure_connections = true
		url = "https://example.com/"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
	}`, name)
}

func secureNotificationChannelWebhookSharedWithCustomData(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_notification_channel_webhook" "sample-webhook5" {
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
