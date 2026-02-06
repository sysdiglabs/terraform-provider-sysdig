//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSSOGroupMappingSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoGroupMappingSettingsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping_settings.test",
						"no_mapping_strategy",
						"UNAUTHORIZED",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping_settings.test",
						"different_roles_same_team_strategy",
						"UNAUTHORIZED",
					),
				),
			},
			{
				Config: ssoGroupMappingSettingsUpdatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping_settings.test",
						"no_mapping_strategy",
						"DEFAULT_TEAM_DEFAULT_ROLE",
					),
				),
			},
			{
				ResourceName:      "sysdig_sso_group_mapping_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ssoGroupMappingSettingsRedirectConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping_settings.test",
						"no_mapping_strategy",
						"NO_MAPPINGS_ERROR_REDIRECT",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_group_mapping_settings.test",
						"no_mappings_error_redirect_url",
						"https://example.com/error",
					),
				),
			},
		},
	})
}

func ssoGroupMappingSettingsConfig() string {
	return `
resource "sysdig_sso_group_mapping_settings" "test" {
  no_mapping_strategy              = "UNAUTHORIZED"
  different_roles_same_team_strategy = "UNAUTHORIZED"
}
`
}

func ssoGroupMappingSettingsUpdatedConfig() string {
	return `
resource "sysdig_sso_group_mapping_settings" "test" {
  no_mapping_strategy              = "DEFAULT_TEAM_DEFAULT_ROLE"
  different_roles_same_team_strategy = "UNAUTHORIZED"
}
`
}

func ssoGroupMappingSettingsRedirectConfig() string {
	return `
resource "sysdig_sso_group_mapping_settings" "test" {
  no_mapping_strategy              = "NO_MAPPINGS_ERROR_REDIRECT"
  different_roles_same_team_strategy = "UNAUTHORIZED"
  no_mappings_error_redirect_url   = "https://example.com/error"
}
`
}
