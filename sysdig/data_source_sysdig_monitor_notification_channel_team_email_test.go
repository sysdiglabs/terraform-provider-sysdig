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

func TestAccMonitorNotificationChannelTeamEmailDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelTeamEmail(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_team_email.nc_team_email", "id", "sysdig_monitor_notification_channel_team_email.nc_team_email", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_team_email.nc_team_email", "name", "sysdig_monitor_notification_channel_team_email.nc_team_email", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_team_email.nc_team_email", "team_id", "sysdig_monitor_notification_channel_team_email.nc_team_email", "team_id"),
				),
			},
		},
	})
}

func monitorNotificationChannelTeamEmail(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_team" "sample_data" {
		name = "monitor-sample-data-%s"
		entrypoint {
		type = "Explore"
		}
	}
resource "sysdig_monitor_notification_channel_team_email" "nc_team_email" {
	name = "%s"
	team_id = sysdig_monitor_team.sample_data.id
}

data "sysdig_monitor_notification_channel_team_email" "nc_team_email" {
	name = sysdig_monitor_notification_channel_team_email.nc_team_email.name
}
`, name, name)
}
