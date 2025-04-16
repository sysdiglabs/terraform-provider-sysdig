//go:build tf_acc_sysdig_monitor || tf_acc_onprem_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSysdigMonitorTeam(t *testing.T) {
	name := fmt.Sprintf("test-monitor-team-%s", randomText(5))
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorTeamResourceAndDatasource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_monitor_team.test", "name", name),
					resource.TestCheckResourceAttr("data.sysdig_monitor_team.test", "description", "A monitor team"),
					resource.TestCheckResourceAttr("data.sysdig_monitor_team.test", "scope_by", "host"),
					resource.TestCheckResourceAttr("data.sysdig_monitor_team.test", "filter", "container.image.repo = \"sysdig/agent\""),
					resource.TestCheckResourceAttr("data.sysdig_monitor_team.test", "can_use_sysdig_capture", "true"),
					resource.TestCheckResourceAttr("data.sysdig_monitor_team.test", "can_see_infrastructure_events", "true"),
					resource.TestCheckResourceAttr("data.sysdig_monitor_team.test", "can_use_aws_data", "true"),
				),
			},
		},
	})
}

func monitorTeamResourceAndDatasource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample" {
  name        = "%s"
  description        = "A monitor team"
  scope_by           			= "host"
  filter             			= "container.image.repo = \"sysdig/agent\""
  can_use_sysdig_capture 		= true
  can_see_infrastructure_events = true
  can_use_aws_data = true
  
  entrypoint {
	type = "Dashboards"
  }
}

data "sysdig_monitor_team" "test" {
  id = sysdig_monitor_team.sample.id
}
`, name)
}
