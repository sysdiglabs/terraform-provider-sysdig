package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorTeamDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

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
				Config: monitorTeam(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_team.acc_team", "name", "sysdig_monitor_team.acc_team", "name"),
					resource.TestCheckResourceAttrPair("data.sysdig_monitor_team.acc_team", "id", "sysdig_monitor_team.acc_team", "id"),
					resource.TestCheckResourceAttr("data.sysdig_monitor_ream.acc_team", "default_team", "true"),
				),
			},
		},
	})
}

func monitorTeam(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "acc_team" {
	name         = "%s"
        default_team = true
}

data "sysdig_monitor_team "acc_team" {
	name = sysdig_monitor_team.acc_team.name
}
`, name)
}
