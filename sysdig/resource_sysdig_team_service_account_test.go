//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

func TestAccTeamServiceAccount(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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
				Config: teamServiceAccount(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account",
						"name",
						name,
					),
					resource.TestCheckResourceAttr("sysdig_team_service_account.service-account",
						"role",
						"ROLE_TEAM_READ",
					),
				),
			},
		},
	})
}

func teamServiceAccount(name string) string {
	return fmt.Sprintf(`
resource "time_static" "example" {
  rfc3339 = "2099-01-01T00:00:00Z"
}

resource "sysdig_monitor_team" "sample" {
  name      = "sample-%s"

  entrypoint {
	type = "Explore"
  }
}

resource "sysdig_team_service_account" "service-account" {
  name = "%s"
  expiration_date = time_static.example.unix
  team_id = sysdig_monitor_team.sample.id
}
`, name, name)
}
