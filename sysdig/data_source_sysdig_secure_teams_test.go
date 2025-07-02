//go:build tf_acc_sysdig_secure || tf_acc_onprem_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSysdigSecureTeams(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSysdigSecureTeamsConfig(randomText(5)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sysdig_secure_teams.test", "teams.0.id"),
				),
			},
		},
	})
}

func testAccDataSourceSysdigSecureTeamsConfig(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "test" {
  name        = "test-secure-team-%s"
  description = "A test secure team"
}

data "sysdig_secure_teams" "test" {}
`, name)
}
