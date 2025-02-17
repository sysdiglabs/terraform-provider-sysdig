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

func TestAccMonitorNotificationChannelMSTeamsDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelMSTeams(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_msteams.nc_msteams", "id", "sysdig_monitor_notification_channel_msteams.nc_msteams", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_msteams.nc_msteams", "name", "sysdig_monitor_notification_channel_msteams.nc_msteams", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_msteams.nc_msteams", "url", "sysdig_monitor_notification_channel_msteams.nc_msteams", "url"),
				),
			},
		},
	})
}

func monitorNotificationChannelMSTeams(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_msteams" "nc_msteams" {
	name = "%s"
	url = "https://sysdig.webhook.office.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
}

data "sysdig_monitor_notification_channel_msteams" "nc_msteams" {
	name = sysdig_monitor_notification_channel_msteams.nc_msteams.name
}
`, name)
}
