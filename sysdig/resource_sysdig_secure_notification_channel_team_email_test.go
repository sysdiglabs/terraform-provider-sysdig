//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common || tf_acc_ibm_secure || tf_acc_ibm_common || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelTeamEmail(t *testing.T) {
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
				Config: secureNotificationChannelTeamEmailWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_team_email.sample_team_email1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: secureNotificationChannelTeamEmailSharedWithCurrentTeam(rText()),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_team_email.sample_team_email2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelTeamEmailWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
	name = "secure-sample-%s"
	all_zones = "true"
}
resource "sysdig_secure_notification_channel_team_email" "sample_team_email1" {
	name = "Example Channel %s - team email"
	enabled = true
	team_id = sysdig_secure_team.sample.id
	notify_when_ok = true
	notify_when_resolved = true
}`, name, name)
}

func secureNotificationChannelTeamEmailSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
	name = "secure-sample-%s"
	all_zones = "true"
}
resource "sysdig_secure_notification_channel_team_email" "sample_team_email2" {
	name = "Example Channel %s - team email"
	enabled = true
	team_id = sysdig_secure_team.sample.id
	notify_when_ok = true
	notify_when_resolved = true
	share_with_current_team = true
}`, name, name)
}
