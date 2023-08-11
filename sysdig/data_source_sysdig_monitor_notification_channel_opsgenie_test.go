//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorNotificationChannelOpsGenieDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelOpsGenie(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_opsgenie.nc_opsgenie", "id", "sysdig_monitor_notification_channel_opsgenie.nc_opsgenie", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_opsgenie.nc_opsgenie", "name", "sysdig_monitor_notification_channel_opsgenie.nc_opsgenie", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_opsgenie.nc_opsgenie", "api_key", "sysdig_monitor_notification_channel_opsgenie.nc_opsgenie", "api_key"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_opsgenie.nc_opsgenie", "region", "sysdig_monitor_notification_channel_opsgenie.nc_opsgenie", "region"),
				),
			},
		},
	})
}

func monitorNotificationChannelOpsGenie(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_opsgenie" "nc_opsgenie" {
	name = "%s"
	api_key = "2349324-342354353-5324-23"
	region = "EU"
}

data "sysdig_monitor_notification_channel_opsgenie" "nc_opsgenie" {
	name = sysdig_monitor_notification_channel_opsgenie.nc_opsgenie.name
}
`, name)
}
