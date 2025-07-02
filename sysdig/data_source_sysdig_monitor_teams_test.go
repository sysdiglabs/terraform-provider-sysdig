//go:build tf_acc_sysdig_monitor || tf_acc_onprem_monitor || tf_acc_ibm_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSysdigMonitorTeams(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigIBMMonitorAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSysdigMonitorTeamsConfig(randomText(10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sysdig_monitor_teams.test", "teams.0.id"),
				),
			},
		},
	})
}

func testAccDataSourceSysdigMonitorTeamsConfig(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample" {
  name        = "%s"
  description        = "A monitor secure team"
  scope_by           			= "host"
  filter             			= "container.image.repo = \"sysdig/agent\""
  can_use_sysdig_capture 		= true
  can_see_infrastructure_events = true
  
  entrypoint {
	type = "Dashboards"
  }
}

data "sysdig_monitor_teams" "test" {}
`, name)
}
