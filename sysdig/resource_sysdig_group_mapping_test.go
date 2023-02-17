package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func TestAccGroupMappingMonitor(t *testing.T) {
	groupMapping1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	groupMapping2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: groupMappingAllTeams("monitor", groupMapping1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("sysdig_monitor_group_mapping.all_teams"),
						"group_name",
						groupMapping1,
					),
					resource.TestCheckResourceAttr(
						"sysdig_monitor_group_mapping.all_teams",
						"team_map.0.all_teams",
						"true",
					),
				),
			},
			{
				Config: groupMappingUpdateAllTeamsGroupName("monitor", groupMapping1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_group_mapping.all_teams",
						"group_name",
						fmt.Sprintf("%s-updated", groupMapping1),
					),
				),
			},
			{
				Config: groupMappingSingleTeamMonitor(groupMapping2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_monitor_group_mapping.single_team",
						"team_map.0.team_ids.#",
						"1",
					),
				),
			},
			{
				ResourceName:      "sysdig_monitor_group_mapping.single_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGroupMappingSecure(t *testing.T) {
	groupMapping1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	groupMapping2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: groupMappingAllTeams("secure", groupMapping1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("sysdig_secure_group_mapping.all_teams"),
						"group_name",
						groupMapping1,
					),
					resource.TestCheckResourceAttr(
						"sysdig_secure_group_mapping.all_teams",
						"team_map.0.all_teams",
						"true",
					),
				),
			},
			{
				Config: groupMappingUpdateAllTeamsGroupName("secure", groupMapping1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_secure_group_mapping.all_teams",
						"group_name",
						fmt.Sprintf("%s-updated", groupMapping1),
					),
				),
			},
			{
				Config: groupMappingSingleTeamSecure(groupMapping2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_secure_group_mapping.single_team",
						"team_map.0.team_ids.#",
						"1",
					),
				),
			},
			{
				ResourceName:      "sysdig_secure_group_mapping.single_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func groupMappingAllTeams(product, groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_%s_group_mapping" "all_teams" {
  group_name = "%s"
  role = "ROLE_TEAM_STANDARD"
  team_map {
    all_teams = true
  }
}
`, product, groupName)
}

func groupMappingUpdateAllTeamsGroupName(product, groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_%s_group_mapping" "all_teams" {
  group_name = "%s-updated"
  role = "ROLE_TEAM_STANDARD"
  team_map {
    all_teams = true
  }
}
`, product, groupName)
}

func groupMappingSingleTeamMonitor(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "single_team" {
  name      = "%[1]s-team-test"

  entrypoint {
	type = "Explore"
  }
}

resource "sysdig_monitor_group_mapping" "single_team" {
  group_name = "%[1]s"
  role = "ROLE_TEAM_STANDARD"

  team_map {
    team_ids = [sysdig_monitor_team.single_team.id]
  }
}
`, groupName)
}

func groupMappingSingleTeamSecure(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "single_team" {
  name      = "%[1]s-team-test"
}

resource "sysdig_secure_group_mapping" "single_team" {
  group_name = "%[1]s"
  role = "ROLE_TEAM_STANDARD"

  team_map {
    team_ids = [sysdig_secure_team.single_team.id]
  }
}
`, groupName)
}
