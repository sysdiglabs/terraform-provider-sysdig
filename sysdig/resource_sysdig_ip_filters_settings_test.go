//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSysdigIpFiltersSettings_fullLifecycle(t *testing.T) {
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
				Config: configBasic(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_ip_filters_settings.test", "ip_filtering_enabled", "true"),
				),
			},
			{
				// Update resource
				Config: configBasic(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_ip_filters_settings.test", "ip_filtering_enabled", "false"),
				),
			},
		},
	})
}

func configBasic(ipFilteringEnabled bool) string {
	return fmt.Sprintf(`
resource "sysdig_ip_filters_settings" "test" {
  ip_filtering_enabled = %t
}
`, ipFilteringEnabled)
}
