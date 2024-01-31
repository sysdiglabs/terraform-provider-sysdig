//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor || tf_acc_onprem_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelWebhookDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelWebhook(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_webhook.nc_webhook", "id", "sysdig_monitor_notification_channel_webhook.nc_webhook", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_webhook.nc_webhook", "name", "sysdig_monitor_notification_channel_webhook.nc_webhook", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_webhook.nc_webhook", "url", "sysdig_monitor_notification_channel_webhook.nc_webhook", "url"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_webhook.nc_webhook", "allow_insecure_connections", "sysdig_monitor_notification_channel_webhook.nc_webhook", "allow_insecure_connections"),
				),
			},
		},
	})
}

func monitorNotificationChannelWebhook(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_webhook" "nc_webhook" {
	name = "%s"
	url = "https://example.com/"
	allow_insecure_connections = false
}

data "sysdig_monitor_notification_channel_webhook" "nc_webhook" {
	name = sysdig_monitor_notification_channel_webhook.nc_webhook.name
}
`, name)
}
