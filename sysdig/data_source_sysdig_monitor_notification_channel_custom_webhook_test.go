//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelCustomWebhookDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: sysdigOrIBMMonitorPreCheck(t),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorNotificationChannelCustomWebhook(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "id", "sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "name", "sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "url", "sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "url"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "http_method", "sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "http_method"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "template", "sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "template"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "allow_insecure_connections", "sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook", "allow_insecure_connections"),
				),
			},
		},
	})
}

func monitorNotificationChannelCustomWebhook(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_custom_webhook" "nc_custom_webhook" {
	name = "%s"
	url = "https://example.com/"
	http_method = "POST"
	template = "{\n  \"code\": \"incident\",\n  \"alert\": \"{{@alert_name}}\"\n}"
	allow_insecure_connections = true
}

data "sysdig_monitor_notification_channel_custom_webhook" "nc_custom_webhook" {
	name = sysdig_monitor_notification_channel_custom_webhook.nc_custom_webhook.name
}
`, name)
}
