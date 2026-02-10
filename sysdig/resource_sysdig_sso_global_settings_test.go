//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSSOGlobalSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: ssoGlobalSettingsConfig("monitor", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_global_settings.test",
						"product",
						"monitor",
					),
					resource.TestCheckResourceAttr(
						"sysdig_sso_global_settings.test",
						"is_password_login_enabled",
						"true",
					),
				),
			},
			{
				Config: ssoGlobalSettingsConfig("monitor", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"sysdig_sso_global_settings.test",
						"is_password_login_enabled",
						"false",
					),
				),
			},
			{
				ResourceName:      "sysdig_sso_global_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func ssoGlobalSettingsConfig(product string, isPasswordLoginEnabled bool) string {
	return fmt.Sprintf(`
resource "sysdig_sso_global_settings" "test" {
  product                    = "%s"
  is_password_login_enabled  = %t
}
`, product, isPasswordLoginEnabled)
}
