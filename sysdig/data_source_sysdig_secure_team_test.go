//go:build tf_acc_sysdig_secure || tf_acc_onprem_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataSourceSysdigSecureTeam(t *testing.T) {
	name := fmt.Sprintf("test-secure-team-%s", randomText(5))
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureTeamAndDatasource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "name", name),
				),
			},
		},
	})
}

func secureTeamAndDatasource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name        = "TF test-%s"
  description = "A test secure team"
}

data "sysdig_secure_team" "test" {
  id = sysdig_secure_team.sample.id
}
`, name)
}
