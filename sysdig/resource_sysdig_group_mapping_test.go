//go:build sysdig_monitor || sysdig_secure

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

func TestAccGroupMapping(t *testing.T) {
	groupAllTeams := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	groupMonitor := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	groupSecure := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
			secure := os.Getenv("SYSDIG_SECURE_API_TOKEN")
			if monitor == "" || secure == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN and SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: groupMappingAllTeams(groupAllTeams),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping.all_teams",
						"group_name",
						groupAllTeams,
					),
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping.all_teams",
						"team_map.0.all_teams",
						"true",
					),
				),
			},
			{
				Config: groupMappingUpdateAllTeamsGroupName(groupAllTeams),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_group_mapping.all_teams",
						"group_name",
						fmt.Sprintf("%s-updated", groupAllTeams),
					),
				),
			},
			{
				Config: groupMappingMonitorTeam(groupMonitor),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("sysdig_group_mapping.group_monitor_%s", groupMonitor),
						"team_map.0.team_ids.#",
						"1",
					),
				),
			},
			{
				ResourceName:      fmt.Sprintf("sysdig_group_mapping.group_monitor_%s", groupMonitor),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: groupMappingSecureTeam(groupSecure),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						fmt.Sprintf("sysdig_group_mapping.group_secure_%s", groupSecure),
						"team_map.0.team_ids.#",
						"1",
					),
				),
			},
			{
				ResourceName:      fmt.Sprintf("sysdig_group_mapping.group_secure_%s", groupSecure),
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
  system_role = "ROLE_USER"

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
  system_role = "ROLE_USER"

  team_map {
    all_teams = true
  }
}
`, groupName)
}

func groupMappingMonitorTeam(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "team_%[1]s" {
  name      = "%[1]s-team-monitor"

  entrypoint {
	type = "Explore"
  }
}

resource "sysdig_group_mapping" "group_monitor_%[1]s" {
  group_name = "%[1]s"
  role = "ROLE_TEAM_STANDARD"

  team_map {
    team_ids = [sysdig_monitor_team.team_%[1]s.id]
  }
}
`, groupName)
}

func groupMappingSecureTeam(groupName string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "team_%[1]s" {
  name      = "%[1]s-team-secure"
}

resource "sysdig_group_mapping" "group_secure_%[1]s" {
  group_name = "%[1]s"
  role = "ROLE_TEAM_STANDARD"

  team_map {
    team_ids = [sysdig_secure_team.team_%[1]s.id]
  }
}
`, groupName)
}
