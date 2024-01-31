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

func TestAccMonitorNotificationChannelIBMCloudFunction(t *testing.T) {
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
				Config: monitorNotificationChannelIBMCloudFunctionWebAction(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_ibm_function.sample1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelIBMCloudFunctionWebActionWithWishAuthToken(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_ibm_function.sample2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelIBMCloudFunctionWebActionWithWishAuthTokenWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_ibm_function.sample3",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelIBMCloudFunctionWebActionWithCustomData(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_ibm_function.sample4",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: monitorNotificationChannelIBMCloudFunctionCloudFunction(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_notification_channel_ibm_function.sample5",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorNotificationChannelIBMCloudFunctionWebAction(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_ibm_function" "sample1" {
	name = "Example Channel %s - IBM Function"
	ibm_function_type = "WEB_ACTION"
	url = "https://eu-gb.functions.cloud.ibm.com/api/v1/web/namespaces/eeeeeeee-623b-4776-ba35-4065bcbfee7b/actions/hello-world/helloworld?param=true"
	whisk_auth_token = "xxx"
}`, name)
}

func monitorNotificationChannelIBMCloudFunctionWebActionWithWishAuthToken(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_ibm_function" "sample2" {
	name = "Example Channel %s - IBM Function"
	ibm_function_type = "WEB_ACTION"
	url = "https://eu-gb.functions.cloud.ibm.com/api/v1/web/namespaces/eeeeeeee-623b-4776-ba35-4065bcbfee7b/actions/hello-world/helloworld?param=true"
	whisk_auth_token = "xxx"
}`, name)
}

func monitorNotificationChannelIBMCloudFunctionWebActionWithWishAuthTokenWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_ibm_function" "sample3" {
	name = "Example Channel %s - IBM Function"
	ibm_function_type = "WEB_ACTION"
	url = "https://eu-gb.functions.cloud.ibm.com/api/v1/web/namespaces/eeeeeeee-623b-4776-ba35-4065bcbfee7b/actions/hello-world/helloworld?param=true"
	share_with_current_team = true
}`, name)
}

func monitorNotificationChannelIBMCloudFunctionWebActionWithCustomData(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_ibm_function" "sample4" {
	name = "Example Channel %s - IBM Function"
	ibm_function_type = "WEB_ACTION"
	url = "https://eu-gb.functions.cloud.ibm.com/api/v1/web/namespaces/eeeeeeee-623b-4776-ba35-4065bcbfee7b/actions/hello-world/helloworld?param=true"
	custom_data = {
		"data1": "value1"
		"data2": "value2"
	}
}`, name)
}

func monitorNotificationChannelIBMCloudFunctionCloudFunction(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_ibm_function" "sample5" {
	name = "Example Channel %s - IBM Function"
	ibm_function_type = "CLOUD_FUNCTION"
	url = "https://eu-gb.functions.cloud.ibm.com/api/v1/namespaces/13eeeeee-623b-4776-ba35-4065bcbfee7b/actions/hello-world/myaction"
	iam_api_key = "xxx"
}`, name)
}
