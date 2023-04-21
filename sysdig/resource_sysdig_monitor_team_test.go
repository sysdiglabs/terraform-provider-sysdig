//go:build tf_acc_sysdig || tf_acc_monitor

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

func TestAccMonitorTeam(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

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
				Config: monitorTeamMinimumConfiguration(rText()),
			},
			{
				Config: monitorTeamWithName(rText()),
			},
			{
				Config: monitorTeamWithFullConfig(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_team.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func monitorTeamWithFullConfig(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample" {
  name                   		= "sample-%s"
  description        			= "%s"
  scope_by           			= "host"
  filter             			= "container.image.repo = \"sysdig/agent\""
  can_use_sysdig_capture 		= true
  can_see_infrastructure_events = true
  can_use_aws_data 				= true
  
  entrypoint {
	type = "Dashboards"
  }
}`, name, name)
}
