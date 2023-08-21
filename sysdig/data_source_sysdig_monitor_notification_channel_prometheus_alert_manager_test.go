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

func TestAccMonitorNotificationChannelPrometheusAlertManagerDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelPrometheusAlertManager(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_prometheus_alert_manager.nc_prometheus_alert_manager", "id", "sysdig_monitor_notification_channel_prometheus_alert_manager.nc_prometheus_alert_manager", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_prometheus_alert_manager.nc_prometheus_alert_manager", "url", "sysdig_monitor_notification_channel_prometheus_alert_manager.nc_prometheus_alert_manager", "url"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_prometheus_alert_manager.nc_prometheus_alert_manager", "allow_insecure_connections", "sysdig_monitor_notification_channel_prometheus_alert_manager.nc_prometheus_alert_manager", "allow_insecure_connections"),
				),
			},
		},
	})
}

func monitorNotificationChannelPrometheusAlertManager(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_prometheus_alert_manager" "nc_prometheus_alert_manager" {
	name = "%s"
	url = "https://testurl.com/xxx"
	allow_insecure_connections = true
}

data "sysdig_monitor_notification_channel_prometheus_alert_manager" "nc_prometheus_alert_manager" {
	name = sysdig_monitor_notification_channel_prometheus_alert_manager.nc_prometheus_alert_manager.name
}
`, name)
}
