//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSSOGroupMappingAllTeams(t *testing.T) {
	groupName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoGroupMappingAllTeamsConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test",
						"group_name",
						groupName,
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test",
						"is_admin",
						"false",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test",
						"team_map.0.is_for_all_teams",
						"true",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test",
						"weight",
						"10",
					),
				),
			},
			{
				ResourceName:      "sysdig_sso_group_mapping.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSSOGroupMappingUpdate(t *testing.T) {
	groupName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoGroupMappingAllTeamsConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test",
						"group_name",
						groupName,
					),
				),
			},
			{
				Config: ssoGroupMappingAllTeamsUpdatedConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test",
						"group_name",
						fmt.Sprintf("%s-updated", groupName),
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test",
						"is_admin",
						"true",
					),
				),
			},
		},
	})
}

func TestAccSSOGroupMappingCustomRole(t *testing.T) {
	groupName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoGroupMappingCustomRoleConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test_custom",
						"group_name",
						groupName,
					),
					resource.TestCheckResourceAttrSet(
						"sysdig_sso_group_mapping.test_custom",
						"custom_team_role_id",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping.test_custom",
						"team_map.0.is_for_all_teams",
						"true",
					),
				),
			},
		},
	})
}

func ssoGroupMappingAllTeamsConfig(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_group_mapping" "test" {
  group_name         = "%s"
  standard_team_role = "ROLE_TEAM_STANDARD"
  is_admin           = false

  team_map {
    is_for_all_teams = true
  }

  weight = 10
}
`, groupName)
}

func ssoGroupMappingAllTeamsUpdatedConfig(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_sso_group_mapping" "test" {
  group_name         = "%s-updated"
  standard_team_role = "ROLE_TEAM_MANAGER"
  is_admin           = true

  team_map {
    is_for_all_teams = true
  }

  weight = 10
}
`, groupName)
}

func ssoGroupMappingCustomRoleConfig(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_custom_role" "test_role" {
  name = "%[1]s-custom-role"
  description = "Test custom role for SSO group mapping"
  permissions {
    monitor_permissions = ["token.view", "api-token.read"]
  }
}

resource "sysdig_sso_group_mapping" "test_custom" {
  group_name         = "%[1]s"
  custom_team_role_id = sysdig_custom_role.test_role.id

  team_map {
    is_for_all_teams = true
  }
}
`, groupName)
}
