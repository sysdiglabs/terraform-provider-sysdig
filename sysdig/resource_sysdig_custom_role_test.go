//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

func TestAccCustomRoleResource(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: customRoleTokenToken(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("sysdig_custom_role.custom-role", "permissions.0.monitor_permissions.*", "token.view"),
					resource.TestCheckTypeSetElemAttr("sysdig_custom_role.custom-role", "permissions.0.monitor_permissions.*", "api-token.read"),
					resource.TestCheckResourceAttr("sysdig_custom_role.custom-role", "permissions.0.secure_permissions.#", "0"),
				),
			},
			{
				Config: customRolePermissionsUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("sysdig_custom_role.custom-role", "permissions.0.monitor_permissions.*", "kubernetes-api-commands.read"),
					resource.TestCheckResourceAttr("sysdig_custom_role.custom-role", "permissions.0.monitor_permissions.#", "1"),
				),
			},
			{
				Config: customRolePermissionsAddSecure(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("sysdig_custom_role.custom-role", "permissions.0.secure_permissions.*", "scanning.read"),
					resource.TestCheckResourceAttr("sysdig_custom_role.custom-role", "permissions.0.secure_permissions.#", "1"),
					resource.TestCheckResourceAttr("sysdig_custom_role.custom-role", "permissions.0.monitor_permissions.#", "1"),
				),
			},
			{
				Config: customRolePermissionsEditName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_custom_role.custom-role",
						"name",
						fmt.Sprintf("custom-role-%s-updated", name),
					),
				),
			},
			{
				ResourceName:      "sysdig_custom_role.custom-role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func customRoleTokenToken(name string) string {
	return fmt.Sprintf(`
resource "sysdig_custom_role" "custom-role" {
  name = "custom-role-%s"
  description = "test"

  permissions {
    monitor_permissions = ["token.view", "api-token.read"]
  }
}
`, name)

}

func customRolePermissionsUpdate(name string) string {
	return fmt.Sprintf(`
resource "sysdig_custom_role" "custom-role" {
  name = "custom-role-%s"
  description = "test"

  permissions {
    monitor_permissions = ["kubernetes-api-commands.read"]
  }
}
`, name)
}

func customRolePermissionsAddSecure(name string) string {
	return fmt.Sprintf(`
resource "sysdig_custom_role" "custom-role" {
  name = "custom-role-%s"
  description = "test"

  permissions {
    monitor_permissions = ["kubernetes-api-commands.read"]
    secure_permissions = ["scanning.read"]
  }
}
`, name)
}

func customRolePermissionsEditName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_custom_role" "custom-role" {
  name = "custom-role-%s-updated"
  description = "test"

  permissions {
    monitor_permissions = ["kubernetes-api-commands.read"]
    secure_permissions = ["scanning.read"]
  }
}
`, name)
}
