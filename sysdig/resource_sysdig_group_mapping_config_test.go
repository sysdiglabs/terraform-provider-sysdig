//go:build tf_acc_sysdig_monitor

package sysdig_test

import (
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

func TestAccGroupMappingConfig(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: groupMappingConfigDefault(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping_config.default",
						"no_mapping_strategy",
						"UNAUTHORIZED",
					),
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping_config.default",
						"different_team_same_role_strategy",
						"UNAUTHORIZED",
					),
				),
			},
			{
				Config: groupMappingConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping_config.default",
						"no_mapping_strategy",
						"DEFAULT_TEAM_DEFAULT_ROLE",
					),
				),
			},
			{
				ResourceName:      "sysdig_group_mapping_config.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func groupMappingConfigDefault() string {
	return `
resource "sysdig_group_mapping_config" "default" {
  no_mapping_strategy = "UNAUTHORIZED"
  different_team_same_role_strategy = "UNAUTHORIZED"
}
`
}

func groupMappingConfigUpdate() string {
	return `
resource "sysdig_group_mapping_config" "default" {
  no_mapping_strategy = "DEFAULT_TEAM_DEFAULT_ROLE"
  different_team_same_role_strategy = "UNAUTHORIZED"
}
`
}
