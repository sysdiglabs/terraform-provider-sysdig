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

func TestAccMonitorNotificationChannelPagerdutyDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelPagerduty(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "name", "sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "account", "sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "account"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "service_key", "sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "service_key"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "service_name", "sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "service_name"),
					resource.TestCheckResourceAttr("data.sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "account", "account"),
					resource.TestCheckResourceAttr("data.sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "service_key", "XXXXXXXXXX"),
					resource.TestCheckResourceAttr("data.sysdig_monitor_notification_channel_pagerduty.nc_pagerduty", "service_name", "sysdig"),
				),
			},
		},
	})
}

func monitorNotificationChannelPagerduty(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_pagerduty" "nc_pagerduty" {
	name = "%s"
	account = "account"
	service_key = "XXXXXXXXXX"
	service_name = "sysdig"
}

data "sysdig_monitor_notification_channel_pagerduty" "nc_pagerduty" {
	name = sysdig_monitor_notification_channel_pagerduty.nc_pagerduty.name
}
`, name)
}
