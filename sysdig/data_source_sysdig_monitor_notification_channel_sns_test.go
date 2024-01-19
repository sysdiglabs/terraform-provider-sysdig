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

func TestAccMonitorNotificationChannelSNSDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelSNS(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_sns.nc_sns", "id", "sysdig_monitor_notification_channel_sns.nc_sns", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_sns.nc_sns", "name", "sysdig_monitor_notification_channel_sns.nc_sns", "name"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_notification_channel_sns.nc_sns", "topics.*", "arn:aws:sns:us-east-1:273489009834:my-alerts2"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_monitor_notification_channel_sns.nc_sns", "topics.*", "arn:aws:sns:us-east-1:279948934544:my-alerts"),
				),
			},
		},
	})
}

func monitorNotificationChannelSNS(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_sns" "nc_sns" {
	name = "%s"
	topics = ["arn:aws:sns:us-east-1:273489009834:my-alerts2", "arn:aws:sns:us-east-1:279948934544:my-alerts"]
}

data "sysdig_monitor_notification_channel_sns" "nc_sns" {
	name = sysdig_monitor_notification_channel_sns.nc_sns.name
}
`, name)
}
