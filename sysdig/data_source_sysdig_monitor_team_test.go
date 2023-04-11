//go:build tf_acc_sysdig || tf_acc_ibm

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

func TestAccMonitorTeamDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: sysdigOrIBMMonitorPreCheck(t),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: monitorTeamDataSourceWithName(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.sysdig_monitor_team.monitor_by_name", "name",
						"sysdig_monitor_team.sample", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.sysdig_monitor_team.monitor_by_name", "entrypoint.0.type",
						"sysdig_monitor_team.sample", "entrypoint.0.type",
					),
				),
			},
		},
	})
}

func monitorTeamDataSourceWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample" {
  name = "sample-%s"
  entrypoint {
	type = "Dashboards"
  }
}

data "sysdig_monitor_team" "monitor_by_name"{
  name = sysdig_monitor_team.sample.name
}
`, name)
}
