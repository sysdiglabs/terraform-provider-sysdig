package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

func TestAccGroupMapping(t *testing.T) {
	groupMapping1 := randomText()
	groupMapping2 := randomText()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckMonitorToken(t),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: groupMappingAllTeams(groupMapping1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping.all_teams",
						"group_name",
						groupMapping1,
					),
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping.all_teams",
						"team_map.0.all_teams",
						"true",
					),
				),
			},
			{
				Config: groupMappingUpdateAllTeamsGroupName(groupMapping1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping.all_teams",
						"group_name",
						fmt.Sprintf("%s-updated", groupMapping1),
					),
				),
			},
			{
				Config: groupMappingSingleTeam(groupMapping2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping.single_team",
						"team_map.0.team_ids.#",
						"1",
					),
				),
			},
			{
				ResourceName:      "sysdig_group_mapping.single_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func groupMappingAllTeams(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_group_mapping" "all_teams" {
  group_name = "%s"
  role = "ROLE_TEAM_STANDARD"
  team_map {
    all_teams = true
  }
}
`, groupName)
}

func groupMappingUpdateAllTeamsGroupName(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_group_mapping" "all_teams" {
  group_name = "%s-updated"
  role = "ROLE_TEAM_STANDARD"
  team_map {
    all_teams = true
  }
}
`, groupName)
}

func groupMappingSingleTeam(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "single_team" {
  name      = "%s-team-test"

  entrypoint {
	type = "Explore"
  }
}

resource "sysdig_group_mapping" "single_team" {
  group_name = "%s"
  role = "ROLE_TEAM_STANDARD"

  team_map {
    team_ids = [sysdig_monitor_team.single_team.id]
  }
}
`, groupName, groupName)
}
