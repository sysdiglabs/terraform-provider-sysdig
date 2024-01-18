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

func TestAccMonitorNotificationChannelVictorOpsDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelVictorOps(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_victorops.nc_victorops", "name", "sysdig_monitor_notification_channel_victorops.nc_victorops", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_victorops.nc_victorops", "account", "sysdig_monitor_notification_channel_victorops.nc_victorops", "account"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_victorops.nc_victorops", "api_key", "sysdig_monitor_notification_channel_victorops.nc_victorops", "api_key"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_victorops.nc_victorops", "routing_key", "sysdig_monitor_notification_channel_victorops.nc_victorops", "routing_key"),
				),
			},
		},
	})
}

func monitorNotificationChannelVictorOps(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_victorops" "nc_victorops" {
	name = "%s"
	api_key = "1234342-4234243-4234-2"
	routing_key = "My team"
}

data "sysdig_monitor_notification_channel_victorops" "nc_victorops" {
	name = sysdig_monitor_notification_channel_victorops.nc_victorops.name
}
`, name)
}
