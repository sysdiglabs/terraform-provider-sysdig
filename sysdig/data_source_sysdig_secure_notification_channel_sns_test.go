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

func TestAccSecureNotificationChannelSNSDataSource(t *testing.T) {
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
				Config: secureNotificationChannelSNS(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_sns.nc_sns", "id", "sysdig_secure_notification_channel_sns.nc_sns", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_sns.nc_sns", "name", "sysdig_secure_notification_channel_sns.nc_sns", "name"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_secure_notification_channel_sns.nc_sns", "topics.*", "arn:aws:sns:us-east-1:273489009834:my-alerts2"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_secure_notification_channel_sns.nc_sns", "topics.*", "arn:aws:sns:us-east-1:279948934544:my-alerts"),
				),
			},
		},
	})
}

func secureNotificationChannelSNS(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_sns" "nc_sns" {
	name = "%s"
	topics = ["arn:aws:sns:us-east-1:273489009834:my-alerts2", "arn:aws:sns:us-east-1:279948934544:my-alerts"]
}

data "sysdig_secure_notification_channel_sns" "nc_sns" {
	name = sysdig_secure_notification_channel_sns.nc_sns.name
}
`, name)
}
