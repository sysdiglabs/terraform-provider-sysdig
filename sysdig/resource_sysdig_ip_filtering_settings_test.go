//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSysdigIpFilteringSettings_fullLifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				// Create resource
				Config: createIPFilteringSettings(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_ip_filtering_settings.test", "ip_filtering_enabled", "false"),
				),
			},
		},
	})
}

func createIPFilteringSettings(ipFilteringEnabled bool) string {
	return fmt.Sprintf(`
resource "sysdig_ip_filtering_settings" "test" {
  ip_filtering_enabled = %t
}
`, ipFilteringEnabled)
}
