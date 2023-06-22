//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common || tf_acc_policies

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelEmail(t *testing.T) {
	t.Cleanup(func() {
		handleSlackNotification(t)
	})

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
				Config: secureNotificationChannelEmailWithName(rText()),
			},
			{
				Config: secureNotificationChannelEmailWithNameInReverseOrder(rText()),
			},
			{
				Config: secureNotificationChannelEmailSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_email.sample_email",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelEmailWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_email" "sample_email" {
	name = "%s"
	recipients = ["root@localhost.com", "bar@localhost.com"]
	enabled = true
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func secureNotificationChannelEmailWithNameInReverseOrder(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_email" "sample_email" {
	name = "%s"
	recipients = ["bar@localhost.com", "root@localhost.com"]
	enabled = false
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func secureNotificationChannelEmailSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_email" "sample_email" {
	name = "%s"
    share_with_current_team = true
	recipients = ["bar@localhost.com", "root@localhost.com"]
	enabled = false
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
