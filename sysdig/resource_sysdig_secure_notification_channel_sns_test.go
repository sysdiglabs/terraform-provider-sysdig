//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelSNS(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelAmazonSNSWithName(rText()),
			},
			{
				Config: secureNotificationChannelAmazonSNSShareWithCurrentTeam(rText()),
			},
			{
				ResourceName:            "sysdig_secure_notification_channel_sns.sample-amazon-sns",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"send_test_notification"},
			},
		},
	})
}

func secureNotificationChannelAmazonSNSWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_sns" "sample-amazon-sns" {
	name = "Example Channel %s - Amazon SNS"
	enabled = true
	topics = ["arn:aws:sns:us-east-1:273489009834:my-alerts2", "arn:aws:sns:us-east-1:279948934544:my-alerts"]
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func secureNotificationChannelAmazonSNSShareWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_sns" "sample-amazon-sns" {
	name = "Example Channel %s - Amazon SNS"
	share_with_current_team = true
	enabled = true
	topics = ["arn:aws:sns:us-east-1:273489009834:my-alerts2", "arn:aws:sns:us-east-1:279948934544:my-alerts"]
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
