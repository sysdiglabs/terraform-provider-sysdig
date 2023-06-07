//go:build tf_acc_sysdig_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccMonitorTeam2(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: createTeamResource(rText()),
			},
			{
				Config: updateTeamResource(rText()),
			},
		},
	})
}

func createTeamResource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team2" "sample" {
  name = "my-team-%s"

}`, name)
}

func updateTeamResource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team2" "sample" {
  name = "updated-my-team-%s"

}`, name)
}
