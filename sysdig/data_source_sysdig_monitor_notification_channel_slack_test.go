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

func TestAccMonitorNotificationChannelSlackDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelSlack(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "id", "sysdig_monitor_notification_channel_slack.nc_slack", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "name", "sysdig_monitor_notification_channel_slack.nc_slack", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "url", "sysdig_monitor_notification_channel_slack.nc_slack", "url"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "show_section_runbook_links", "sysdig_monitor_notification_channel_slack.nc_slack", "show_section_runbook_links"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "show_section_event_details", "sysdig_monitor_notification_channel_slack.nc_slack", "show_section_event_details"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "show_section_user_defined_content", "sysdig_monitor_notification_channel_slack.nc_slack", "show_section_user_defined_content"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "show_section_notification_chart", "sysdig_monitor_notification_channel_slack.nc_slack", "show_section_notification_chart"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "show_section_dashboard_links", "sysdig_monitor_notification_channel_slack.nc_slack", "show_section_dashboard_links"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "show_section_alert_details", "sysdig_monitor_notification_channel_slack.nc_slack", "show_section_alert_details"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_slack.nc_slack", "show_section_capturing_information", "sysdig_monitor_notification_channel_slack.nc_slack", "show_section_capturing_information"),
				),
			},
		},
	})
}

func monitorNotificationChannelSlack(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_slack" "nc_slack" {
	name = "%s"
	url = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	show_section_runbook_links = false
	show_section_event_details = false
	show_section_user_defined_content = false
	show_section_notification_chart = false
	show_section_dashboard_links = false
	show_section_alert_details = false
	show_section_capturing_information = false
}

data "sysdig_monitor_notification_channel_slack" "nc_slack" {
	name = sysdig_monitor_notification_channel_slack.nc_slack.name
}
`, name)
}
