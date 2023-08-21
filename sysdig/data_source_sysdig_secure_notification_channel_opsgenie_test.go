//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelOpsGenieDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelOpsGenie(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_opsgenie.nc_opsgenie", "id", "sysdig_secure_notification_channel_opsgenie.nc_opsgenie", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_opsgenie.nc_opsgenie", "name", "sysdig_secure_notification_channel_opsgenie.nc_opsgenie", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_opsgenie.nc_opsgenie", "api_key", "sysdig_secure_notification_channel_opsgenie.nc_opsgenie", "api_key"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_opsgenie.nc_opsgenie", "region", "sysdig_secure_notification_channel_opsgenie.nc_opsgenie", "region"),
				),
			},
		},
	})
}

func secureNotificationChannelOpsGenie(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_opsgenie" "nc_opsgenie" {
	name = "%s"
	api_key = "2349324-342354353-5324-23"
	region = "EU"
}

data "sysdig_secure_notification_channel_opsgenie" "nc_opsgenie" {
	name = sysdig_secure_notification_channel_opsgenie.nc_opsgenie.name
}
`, name)
}
