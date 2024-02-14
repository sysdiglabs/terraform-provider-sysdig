//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccTeamServiceAccount(t *testing.T) {
	monitorsvc := randomText(10)
	securesvc := randomText(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {Source: "hashicorp/time"},
		},
		Steps: []resource.TestStep{
			{
				Config: teamServiceAccountMonitorTeam(monitorsvc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account-monitor",
						"name",
						monitorsvc,
					),
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account-monitor",
						"role",
						"ROLE_TEAM_READ",
					),
					resource.TestCheckResourceAttrSet("sysdig_team_service_account.service-account-monitor",
						"api_key",
					),
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account-monitor",
						"expiration_date",
						"4070908800",
					),
				),
			},
			{
				Config: teamServiceAccountMonitorTeamNewExpirationDate(monitorsvc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account-monitor",
						"name",
						monitorsvc,
					),
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account-monitor",
						"role",
						"ROLE_TEAM_READ",
					),
					resource.TestCheckResourceAttrSet("sysdig_team_service_account.service-account-monitor",
						"api_key",
					),
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account-monitor",
						"expiration_date",
						"4070995200",
					),
				),
			},
			{
				Config: teamServiceAccountSecureTeam(securesvc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account-secure",
						"name",
						securesvc,
					),
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account-secure",
						"role",
						"ROLE_TEAM_READ",
					),
					resource.TestCheckResourceAttrSet("sysdig_team_service_account.service-account-secure",
						"api_key",
					),
				),
			},
		},
	})
}

func teamServiceAccountMonitorTeam(name string) string {
	return fmt.Sprintf(`
resource "time_static" "example" {
  rfc3339 = "2099-01-01T00:00:00Z"
}

resource "sysdig_monitor_team" "sample" {
  name      = "monitor-sample-%s"

  entrypoint {
	type = "Explore"
  }
}

resource "sysdig_team_service_account" "service-account-monitor" {
  name = "%s"
  expiration_date = time_static.example.unix
  team_id = sysdig_monitor_team.sample.id
}
`, name, name)
}

func teamServiceAccountMonitorTeamNewExpirationDate(name string) string {
	return fmt.Sprintf(`
resource "time_static" "example" {
  rfc3339 = "2099-01-02T00:00:00Z"
}

resource "sysdig_monitor_team" "sample" {
  name      = "monitor-sample-%s"

  entrypoint {
	type = "Explore"
  }
}

resource "sysdig_team_service_account" "service-account-monitor" {
  name = "%s"
  expiration_date = time_static.example.unix
  team_id = sysdig_monitor_team.sample.id
}
`, name, name)
}

func teamServiceAccountSecureTeam(name string) string {
	return fmt.Sprintf(`
resource "time_static" "example" {
  rfc3339 = "2099-01-01T00:00:00Z"
}

resource "sysdig_secure_team" "sample" {
  name      = "secure-sample-%s"
  all_zones = "true"
}

resource "sysdig_team_service_account" "service-account-secure" {
  name = "%s"
  expiration_date = time_static.example.unix
  team_id = sysdig_secure_team.sample.id
}
`, name, name)
}
