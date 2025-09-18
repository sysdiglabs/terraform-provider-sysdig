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

func TestAccMonitorSilenceRule(t *testing.T) {
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
				Config: monitorSilenceRuleWithName(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_silence_rule.sample1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorSilenceRuleWithAlertIds(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_silence_rule.sample2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorSilenceRuleWithAlertIdsAndScope(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_silence_rule.sample3",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorSilenceRuleWithNotificationChannels(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_silence_rule.sample4",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorSilenceRuleWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_silence_rule" "sample1" {
	name = "Example Silence Rule %s"
	enabled = false
	start_ts = 1691168134153
	duration_seconds = 3600
	scope = "container.name in (\"test\")"
}`, name)
}

func monitorSilenceRuleWithAlertIds(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_prometheus" "sample1" {
	name = "TERRAFORM TEST - PROMQL %s 1"
	query = "up"
	duration_seconds = 60
	enabled = false
}
resource "sysdig_monitor_alert_v2_prometheus" "sample2" {
	name = "TERRAFORM TEST - PROMQL %s 2"
	query = "up"
	duration_seconds = 60
	enabled = false
}
resource "sysdig_monitor_silence_rule" "sample2" {
	name = "Example Silence Rule %s"
	enabled = false
	start_ts = 1691168134153
	duration_seconds = 3600
	alert_ids = [ sysdig_monitor_alert_v2_prometheus.sample1.id, sysdig_monitor_alert_v2_prometheus.sample2.id ]
}`, name, name, name)
}

func monitorSilenceRuleWithAlertIdsAndScope(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_prometheus" "sample3" {
	name = "TERRAFORM TEST - PROMQL %s 3"
	query = "up"
	duration_seconds = 60
	enabled = false
}
resource "sysdig_monitor_alert_v2_prometheus" "sample4" {
	name = "TERRAFORM TEST - PROMQL %s 4"
	query = "up"
	duration_seconds = 60
	enabled = false
}
resource "sysdig_monitor_silence_rule" "sample3" {
	name = "Example Silence Rule %s"
	enabled = false
	start_ts = 1691168134153
	duration_seconds = 3600
	scope = "container.name in (\"test\")"
	alert_ids = [ sysdig_monitor_alert_v2_prometheus.sample3.id, sysdig_monitor_alert_v2_prometheus.sample4.id ]
}`, name, name, name)
}

func monitorSilenceRuleWithNotificationChannels(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_webhook" "sample-webhook" {
	name = "Example Channel %s - Webhook"
	enabled = false
	url = "https://example.com/"
	send_test_notification = false
}
resource "sysdig_monitor_silence_rule" "sample4" {
	name = "Example Silence Rule %s"
	enabled = false
	start_ts = 1691168134153
	duration_seconds = 3600
	scope = "container.name in (\"test\")"
	notification_channel_ids = [sysdig_monitor_notification_channel_webhook.sample-webhook.id]
}`, name, name)
}
