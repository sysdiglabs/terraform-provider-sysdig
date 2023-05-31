//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelSlack(t *testing.T) {
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
				Config: secureNotificationChannelSlackWithName(rText()),
			},
			{
				Config: secureNotificationChannelSlackSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_slack.sample-slack",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelSlackWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_slack" "sample-slack" {
	name = "Example Channel %s - Slack"
	enabled = true
	url = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}

func secureNotificationChannelSlackSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_slack" "sample-slack" {
	name = "Example Channel %s - Slack"
    share_with_current_team = true
	enabled = true
	url = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel = "#sysdig"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}
