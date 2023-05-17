//go:build tf_acc_sysdig || tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelVictorOps(t *testing.T) {
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
				Config: secureNotificationChannelVictorOpsWithName(rText()),
			},
			{
				Config: secureNotificationChannelVictorOpsShareWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_victorops.sample-victorops",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelVictorOpsWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_victorops" "sample-victorops" {
	name = "Example Channel %s - VictorOps"
	enabled = true
	api_key = "1234342-4234243-4234-2"
	routing_key = "My team"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}

func secureNotificationChannelVictorOpsShareWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_victorops" "sample-victorops" {
	name = "Example Channel %s - VictorOps"
    share_with_current_team = true
	enabled = true
	api_key = "1234342-4234243-4234-2"
	routing_key = "My team"
	notify_when_ok = false
	notify_when_resolved = false
	send_test_notification = false
}`, name)
}
