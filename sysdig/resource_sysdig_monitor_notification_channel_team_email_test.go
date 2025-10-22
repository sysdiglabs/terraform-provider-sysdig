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

func TestAccMonitorNotificationChannelTeamEmail(t *testing.T) {
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
				Config: monitorNotificationChannelTeamEmailWithName(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_team_email.sample_team_email1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: monitorNotificationChannelTeamEmailSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_team_email.sample_team_email2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
			{
				Config: monitorNotificationChannelTeamEmailWithIncludeAdminUsers(rText()),
			},
			{
				ResourceName:            "sysdig_monitor_notification_channel_team_email.sample_team_email3",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
		},
	})
}

func monitorNotificationChannelTeamEmailWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample1" {
	name = "monitor-sample-%s"
	entrypoint {
		type = "Explore"
	}
}
resource "sysdig_monitor_notification_channel_team_email" "sample_team_email1" {
	name = "Example Channel %s - team email"
	enabled = true
	team_id = sysdig_monitor_team.sample1.id
	notify_when_ok = true
	notify_when_resolved = true
}`, name, name)
}

func monitorNotificationChannelTeamEmailSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample2" {
	name = "monitor-sample-%s"
	entrypoint {
		type = "Explore"
	}
}
resource "sysdig_monitor_notification_channel_team_email" "sample_team_email2" {
	name = "Example Channel %s - team email"
	enabled = true
	team_id = sysdig_monitor_team.sample2.id
	notify_when_ok = true
	notify_when_resolved = true
	share_with_current_team = true
}`, name, name)
}

func monitorNotificationChannelTeamEmailWithIncludeAdminUsers(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample3" {
	name = "monitor-sample-%s"
	entrypoint {
		type = "Explore"
	}
}
resource "sysdig_monitor_notification_channel_team_email" "sample_team_email3" {
	name = "Example Channel %s - team email"
	enabled = true
	team_id = sysdig_monitor_team.sample3.id
	include_admin_users = true
	notify_when_ok = true
	notify_when_resolved = true
}`, name, name)
}
