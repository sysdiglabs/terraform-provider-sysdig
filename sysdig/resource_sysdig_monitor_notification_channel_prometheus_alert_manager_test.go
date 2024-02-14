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

func TestAccMonitorNotificationChannelPrometheusAlertManager(t *testing.T) {
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
				Config: monitorNotificationChannelPrometheusAlertManagerWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_prometheus_alert_manager.sample-channel1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelPrometheusAlertManagerWithNameWithAdditionalheaders(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_prometheus_alert_manager.sample-channel2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelPrometheusAlertManagerWithNameWithAllowInsecureConnections(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_prometheus_alert_manager.sample-channel3",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelPrometheusAlertManagerSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_prometheus_alert_manager.sample-channel4",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelPrometheusAlertManagerWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_prometheus_alert_manager" "sample-channel1" {
	name = "Example Channel %s - prometheus alert manager"
	enabled = true
	url = "https://testurl.com/xxx"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func monitorNotificationChannelPrometheusAlertManagerWithNameWithAdditionalheaders(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_prometheus_alert_manager" "sample-channel2" {
		name = "Example Channel %s - prometheus alert manager With Additional Headers"
		enabled = true
		url = "https://testurl.com/xxx"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		additional_headers = {
			"custom-Header": "TestHeader"
		}
	}`, name)
}

func monitorNotificationChannelPrometheusAlertManagerWithNameWithAllowInsecureConnections(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_prometheus_alert_manager" "sample-channel3" {
		name = "Example Channel %s - prometheus alert manager with insecure connections"
		enabled = true
		url = "https://testurl.com/xxx"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		allow_insecure_connections = true
	}`, name)
}

func monitorNotificationChannelPrometheusAlertManagerSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_notification_channel_prometheus_alert_manager" "sample-channel4" {
		name = "Example Channel %s - prometheus alert manager with share with current team"
		enabled = true
		url = "https://testurl.com/xxx"
		notify_when_ok = false
		notify_when_resolved = false
		send_test_notification = false
		share_with_current_team = true
	}`, name)
}
