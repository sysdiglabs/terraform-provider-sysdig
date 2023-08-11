//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelSlackDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelSlack(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "id", "sysdig_secure_notification_channel_slack.nc_slack", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "name", "sysdig_secure_notification_channel_slack.nc_slack", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "url", "sysdig_secure_notification_channel_slack.nc_slack", "url"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "show_section_runbook_links", "sysdig_secure_notification_channel_slack.nc_slack", "show_section_runbook_links"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "show_section_event_details", "sysdig_secure_notification_channel_slack.nc_slack", "show_section_event_details"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "show_section_user_defined_content", "sysdig_secure_notification_channel_slack.nc_slack", "show_section_user_defined_content"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "show_section_notification_chart", "sysdig_secure_notification_channel_slack.nc_slack", "show_section_notification_chart"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "show_section_notification_chart", "sysdig_secure_notification_channel_slack.show_section_dashboard_links", "show_section_dashboard_links"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "show_section_notification_chart", "sysdig_secure_notification_channel_slack.show_section_alert_details", "show_section_alert_details"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_slack.nc_slack", "show_section_notification_chart", "sysdig_secure_notification_channel_slack.show_section_capturing_information", "show_section_capturing_information"),
				),
			},
		},
	})
}

func secureNotificationChannelSlack(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_slack" "nc_slack" {
	name = "%s"
	url = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	template_version = "v2"
}

data "sysdig_secure_notification_channel_slack" "nc_slack" {
	name = sysdig_secure_notification_channel_slack.nc_slack.name
}
`, name)
}
