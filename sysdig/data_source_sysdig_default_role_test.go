//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccDataSourceSysdigDefaultRole(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `data "sysdig_default_role" "advanced" {
  name = "Advanced User"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_default_role.advanced", "name", "Advanced User"),
					// Verify both permission sets are non-empty
					resource.TestCheckResourceAttrSet("data.sysdig_default_role.advanced", "monitor_permissions.#"),
					resource.TestCheckResourceAttrSet("data.sysdig_default_role.advanced", "secure_permissions.#"),
					// Verify well-known monitor permissions are present
					resource.TestCheckTypeSetElemAttr("data.sysdig_default_role.advanced", "monitor_permissions.*", "alerts.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_default_role.advanced", "monitor_permissions.*", "dashboards.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_default_role.advanced", "monitor_permissions.*", "token.view"),
					// Verify well-known secure permissions are present
					resource.TestCheckTypeSetElemAttr("data.sysdig_default_role.advanced", "secure_permissions.*", "scanning.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_default_role.advanced", "secure_permissions.*", "secure.policy.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_default_role.advanced", "secure_permissions.*", "policies.read"),
				),
			},
		},
	})
}
