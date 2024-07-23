//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSysdigAllowedIpRange_fullLifecycle(t *testing.T) {
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
				Config: configBasic("192.168.1.0/24", "Initial note", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_allowed_ip_range.test", "ip_range", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("sysdig_allowed_ip_range.test", "note", "Initial note"),
					resource.TestCheckResourceAttr("sysdig_allowed_ip_range.test", "enabled", "true"),
				),
			},
			{
				// Update resource
				Config: configBasic("192.168.2.0/24", "Updated note", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_allowed_ip_range.test", "ip_range", "192.168.2.0/24"),
					resource.TestCheckResourceAttr("sysdig_allowed_ip_range.test", "note", "Updated note"),
					resource.TestCheckResourceAttr("sysdig_allowed_ip_range.test", "enabled", "false"),
				),
			},
		},
	})
}

func configBasic(ipRange, note string, enabled bool) string {
	return fmt.Sprintf(`
resource "sysdig_allowed_ip_range" "test" {
  ip_range = "%s"
  note     = "%s"
  enabled  = %t
}
`, ipRange, note, enabled)
}
