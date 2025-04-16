//go:build tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccDataSourceSysdigSecureTeamIBM(t *testing.T) {
	name := fmt.Sprintf("test-secure-team-%s", randomText(5))
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureTeamResourceAndDatasource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test_dt", "name", name),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test_dt", "description", "A secure team"),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test_dt", "enable_ibm_platform_metrics", "true"),
					resource.TestCheckResourceAttr("data.sysdig_secure_team.test_dt", "ibm_platform_metrics", "foo in (\"0\") and bar in (\"3\")"),
				),
			},
		},
	})
}

func secureTeamWithPlatformMetricsIBM(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "test" {
  name = "%s"
  enable_ibm_platform_metrics = true
  ibm_platform_metrics = "foo in (\"0\") and bar in (\"3\")"

  entrypoint {
	type = "Dashboards"
  }
}

data "sysdig_secure_team" "test_dt" {
  id = sysdig_secure_team.sample.id
}
`, name)
}
