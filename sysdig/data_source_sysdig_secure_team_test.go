//go:build tf_acc_sysdig_secure || tf_acc_onprem_secure

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
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "description", "A test secure team"),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "scope_by", "container"),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "filter", "container.image.repo = \"sysdig/agent\""),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "version", "0"),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "can_use_sysdig_capture", "true"),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "can_see_infrastructure_events", "true"),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "can_use_aws_data", "false"),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test", "all_zones", "true"),
				),
			},
		},
	})
}

func secureTeamAndDatasource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name               = "%s"
  description        = "A test secure team"
  scope_by           = "container"
  can_use_sysdig_capture = true
  can_see_infrastructure_events = true
  filter             = "container.image.repo = \"sysdig/agent\""
  can_use_aws_data   = false
  all_zones          = true
}

data "sysdig_secure_team" "test" {
  id = sysdig_secure_team.sample.id
}
`, name)
}
