//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccSysdigIpFilter_fullLifecycle(t *testing.T) {
	ipRange1 := generateRandomIPRange()
	ipRange2 := generateRandomIPRange()

	resource.Test(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				// Create resource with the first random IP range
				Config: createIPFilter(ipRange1, "Initial note", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_ip_filter.test", "ip_range", ipRange1),
					resource.TestCheckResourceAttr("sysdig_ip_filter.test", "note", "Initial note"),
					resource.TestCheckResourceAttr("sysdig_ip_filter.test", "enabled", "true"),
				),
			},
			{
				// Update resource with the second random IP range
				Config: createIPFilter(ipRange2, "Updated note", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_ip_filter.test", "ip_range", ipRange2),
					resource.TestCheckResourceAttr("sysdig_ip_filter.test", "note", "Updated note"),
					resource.TestCheckResourceAttr("sysdig_ip_filter.test", "enabled", "false"),
				),
			},
		},
	})
}

func generateRandomIPRange() string {
	return fmt.Sprintf("%d.%d.%d.0/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(0, 255), acctest.RandIntRange(0, 255))
}

func createIPFilter(ipRange, note string, enabled bool) string {
	return fmt.Sprintf(`
resource "sysdig_ip_filter" "test" {
  ip_range = "%s"
  note     = "%s"
  enabled  = %t
}
`, ipRange, note, enabled)
}
