//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccCustomRoleDateSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getCustomRole(rText),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("data.sysdig_custom_role.custom", "monitor_permissions.*", "token.view"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_custom_role.custom", "monitor_permissions.*", "api-token.read"),
					resource.TestCheckResourceAttr("data.sysdig_custom_role.custom", "secure_permissions.#", "0"),
				),
			},
		},
	})
}

func getCustomRole(name string) string {
	return fmt.Sprintf(`
resource "sysdig_custom_role" "test" {
  name = "%s"
  description = "test"

  permissions {
    monitor_permissions = ["token.view", "api-token.read"]
  }
}
data "sysdig_custom_role" "custom" {
  depends_on = [sysdig_custom_role.test]
  name = sysdig_custom_role.test.name
}
`, name)
}
