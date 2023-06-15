//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_common || tf_acc_ibm_monitor || tf_acc_ibm_common

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelMSTeams(t *testing.T) {
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
				Config: monitorNotificationChannelMSTeamsWithName(rText()),
			},
			{
				Config: monitorNotificationChannelMSTeamsSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_msteams.sample-msteams",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelMSTeamsWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_msteams" "sample-msteams" {
	name = "Example Channel %s - MS Teams"
	enabled = true
	url = "https://hooks.msteams.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}

func monitorNotificationChannelMSTeamsSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_msteams" "sample-msteams" {
	name = "Example Channel %s - MS Teams"
    share_with_current_team = true
	enabled = true
	url = "https://hooks.msteams.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}
