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

func TestAccSecureNotificationChannelMSTeamsDataSource(t *testing.T) {
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
				Config: secureNotificationChannelMSTeams(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_msteams.nc_msteams", "id", "sysdig_secure_notification_channel_msteams.nc_msteams", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_msteams.nc_msteams", "name", "sysdig_secure_notification_channel_msteams.nc_msteams", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_msteams.nc_msteams", "url", "sysdig_secure_notification_channel_msteams.nc_msteams", "url"),
				),
			},
		},
	})
}

func secureNotificationChannelMSTeams(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_msteams" "nc_msteams" {
	name = "%s"
	url = "https://hooks.msteams.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	template_version = "v2"
}

data "sysdig_secure_notification_channel_msteams" "nc_msteams" {
	name = sysdig_secure_notification_channel_msteams.nc_msteams.name
}
`, name)
}
