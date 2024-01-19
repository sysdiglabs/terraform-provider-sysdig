//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelVictorOpsDataSource(t *testing.T) {
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
				Config: secureNotificationChannelVictorOps(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_victorops.nc_victorops", "name", "sysdig_secure_notification_channel_victorops.nc_victorops", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_victorops.nc_victorops", "account", "sysdig_secure_notification_channel_victorops.nc_victorops", "account"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_victorops.nc_victorops", "api_key", "sysdig_secure_notification_channel_victorops.nc_victorops", "api_key"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_victorops.nc_victorops", "routing_key", "sysdig_secure_notification_channel_victorops.nc_victorops", "routing_key"),
				),
			},
		},
	})
}

func secureNotificationChannelVictorOps(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_victorops" "nc_victorops" {
	name = "%s"
	api_key = "1234342-4234243-4234-2"
	routing_key = "My team"
}

data "sysdig_secure_notification_channel_victorops" "nc_victorops" {
	name = sysdig_secure_notification_channel_victorops.nc_victorops.name
}
`, name)
}
