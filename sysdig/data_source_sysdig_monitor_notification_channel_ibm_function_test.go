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

func TestAccMonitorNotificationChannelIBMFunctionDataSource(t *testing.T) {
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
				Config: monitorNotificationChannelIBMFunction(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "id", "sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "name", "sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "ibm_function_type", "sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "ibm_function_type"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "url", "sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "url"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "whisk_auth_token", "sysdig_monitor_notification_channel_ibm_function.nc_ibm_function", "whisk_auth_token"),
				),
			},
		},
	})
}

func monitorNotificationChannelIBMFunction(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_ibm_function" "nc_ibm_function" {
	name = "%s"
	ibm_function_type = "WEB_ACTION"
	url = "https://eu-gb.functions.cloud.ibm.com/api/v1/web/namespaces/eeeeeeee-623b-4776-ba35-4065bcbfee7b/actions/hello-world/helloworld?param=true"
	whisk_auth_token = "xxx"
}

data "sysdig_monitor_notification_channel_ibm_function" "nc_ibm_function" {
	name = sysdig_monitor_notification_channel_ibm_function.nc_ibm_function.name
}
`, name)
}
