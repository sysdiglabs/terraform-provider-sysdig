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

func TestAccMonitorNotificationChannelSlack(t *testing.T) {
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
				Config: monitorNotificationChannelSlackWithName(rText()),
			},
			{
				Config: monitorNotificationChannelSlackSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_slack.sample-slack",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelSlackSharedWithShowSection(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_slack.sample-slack",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelSlackSharedWithPrivateChannel(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_slack.sample-slack-private",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelSlackWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_slack" "sample-slack" {
	name = "Example Channel %s - Slack"
	enabled = true
	url = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}

func monitorNotificationChannelSlackSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_slack" "sample-slack" {
	name = "Example Channel %s - Slack"
	share_with_current_team = true
	enabled = true
	url = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}

func monitorNotificationChannelSlackSharedWithShowSection(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_slack" "sample-slack" {
	name = "Example Channel %s - Slack"
	enabled = true
	url = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	notify_when_ok = true
	notify_when_resolved = true
	show_section_runbook_links = false
	show_section_event_details = false
	show_section_user_defined_content = false
	show_section_notification_chart = false
	show_section_dashboard_links = false
	show_section_alert_details = false
	show_section_capturing_information = false
}`, name)
}

func monitorNotificationChannelSlackSharedWithPrivateChannel(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_slack" "sample-slack-private" {
	name = "Example Channel %s - Slack"
	enabled = true
	url = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	private_channel = true
	private_channel_url = "https://app.slack.com/client/XXXXXXXX/XXXXXXXX"
}`, name)
}
