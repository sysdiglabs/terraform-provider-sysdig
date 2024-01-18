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

func TestAccSecureNotificationChannelEmailDataSource(t *testing.T) {
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
				Config: secureNotificationChannelEmail(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_email.nc_email", "id", "sysdig_secure_notification_channel_email.nc_email", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_email.nc_email", "name", "sysdig_secure_notification_channel_email.nc_email", "name"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_secure_notification_channel_email.nc_email", "recipients.*", "root@localhost.com"),
				),
			},
		},
	})
}

func secureNotificationChannelEmail(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_email" "nc_email" {
	name = "%s"
	recipients = ["root@localhost.com"]
}

data "sysdig_secure_notification_channel_email" "nc_email" {
	name = sysdig_secure_notification_channel_email.nc_email.name
}
`, name)
}
