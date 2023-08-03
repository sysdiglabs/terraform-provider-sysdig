//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureNotificationChannelMSTeams(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureNotificationChannelMSTeamsWithName(rText()),
			},
			{
				Config: secureNotificationChannelMSTeamsSharedWithCurrentTeam(rText()),
			},
			{
				Config: secureNotificationChannelMSTeamsWithNameAndTemplateVersion(rText(), "v2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_notification_channel_msteams.sample-msteams", "template_version", "v2"),
				),
			},
			{
				Config: secureNotificationChannelMSTeamsWithNameAndTemplateVersion(rText(), "v1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_notification_channel_msteams.sample-msteams", "template_version", "v1"),
				),
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_msteams.sample-msteams",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureNotificationChannelMSTeamsWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_msteams" "sample-msteams" {
	name = "Example Channel %s - MS Teams"
	enabled = true
	url = "https://hooks.msteams.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}

func secureNotificationChannelMSTeamsSharedWithCurrentTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_msteams" "sample-msteams" {
	name = "Example Channel %s - MS Teams"
	share_with_current_team = true
	enabled = true
	url = "https://hooks.msteams.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	notify_when_ok = true
	notify_when_resolved = true
}`, name)
}

func secureNotificationChannelMSTeamsWithNameAndTemplateVersion(name, version string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_msteams" "sample-msteams" {
	name = "Example Channel %s - MS Teams"
	enabled = true
	url = "https://hooks.msteams.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	notify_when_ok = true
	notify_when_resolved = true
	template_version = "%s"
}`, name, version)
}
