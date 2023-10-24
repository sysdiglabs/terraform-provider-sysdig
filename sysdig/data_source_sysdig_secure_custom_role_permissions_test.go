//go:build tf_acc_sysdig_secure

package sysdig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureCustomRolePermissionsDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureCustomRolePermissions(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("data.sysdig_secure_custom_role_permissions.images_edit", "enriched_permissions.*", "secure.blacklist.images.edit"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_secure_custom_role_permissions.images_edit", "enriched_permissions.*", "secure.blacklist.images.read"),
					resource.TestCheckTypeSetElemAttr("data.sysdig_secure_custom_role_permissions.images_edit", "enriched_permissions.*", "scanning.read"),
					resource.TestCheckResourceAttr("data.sysdig_secure_custom_role_permissions.images_edit", "enriched_permissions.#", "3"),
				),
			},
		},
	})
}

func secureCustomRolePermissions() string {
	return `
data "sysdig_secure_custom_role_permissions" "images_edit" {
  requested_permissions = ["secure.blacklist.images.edit"]
}
`
}
