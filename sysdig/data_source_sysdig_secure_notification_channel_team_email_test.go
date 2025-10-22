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

func TestAccSecureNotificationChannelTeamEmailDataSource(t *testing.T) {
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
				Config: secureNotificationChannelTeamEmail(rText()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_team_email.nc_team_email", "id", "sysdig_secure_notification_channel_team_email.nc_team_email", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_team_email.nc_team_email", "name", "sysdig_secure_notification_channel_team_email.nc_team_email", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_team_email.nc_team_email", "team_id", "sysdig_secure_notification_channel_team_email.nc_team_email", "team_id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_team_email.nc_team_email", "include_admin_users", "sysdig_secure_notification_channel_team_email.nc_team_email", "include_admin_users"),
				),
			},
			{
				Config: resourceSecureNotificationChannelTeamEmailWithIncludeAdminUsers(rText()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_team_email.nc_team_email_2", "id", "sysdig_secure_notification_channel_team_email_2.nc_team_email", "id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_team_email.nc_team_email_2", "name", "sysdig_secure_notification_channel_team_email_2.nc_team_email", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_team_email.nc_team_email_2", "team_id", "sysdig_secure_notification_channel_team_email_2.nc_team_email", "team_id"),
					resource.TestCheckResourceAttrPair("data.sysdig_secure_notification_channel_team_email.nc_team_email_2", "include_admin_users", "sysdig_secure_notification_channel_team_email_2.nc_team_email", "include_admin_users"),
				),
			},
		},
	})
}

func secureNotificationChannelTeamEmail(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_team" "sample_data" {
		name = "secure-sample-data-%s"
		all_zones = "true"
	}
resource "sysdig_secure_notification_channel_team_email" "nc_team_email" {
	name = "%s"
	team_id = sysdig_secure_team.sample_data.id
}

data "sysdig_secure_notification_channel_team_email" "nc_team_email" {
	name = sysdig_secure_notification_channel_team_email.nc_team_email.name
}
`, name, name)
}

func resourceSecureNotificationChannelTeamEmailWithIncludeAdminUsers(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_team" "sample_data_2" {
		name = "secure-sample-data-%s"
		all_zones = "true"
	}
resource "sysdig_secure_notification_channel_team_email" "nc_team_email_2" {
	name = "%s"
	team_id = sysdig_secure_team.sample_data_2.id
	include_admin_users = true
}

data "sysdig_secure_notification_channel_team_email" "nc_team_email_2" {
	name = sysdig_secure_notification_channel_team_email.nc_team_email_2.name
}
`, name, name)
}
