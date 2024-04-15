//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccAgentAccessKeyDataSource(t *testing.T) {
	limit := 1
	reservation := 0
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getAgentAccessKey(limit, reservation, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "limit", strconv.Itoa(limit)),
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "reservation", strconv.Itoa(reservation)),
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "metadata.test", "yes"),
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "metadata.acceptance_test", "true"),
				),
			},
			{
				Config: getAgentAccessKey(limit, reservation, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "limit", strconv.Itoa(limit)),
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "reservation", strconv.Itoa(reservation)),
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "enabled", strconv.FormatBool(false)),
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "metadata.test", "yes"),
					resource.TestCheckResourceAttr("data.sysdig_agent_access_key.data", "metadata.acceptance_test", "true"),
				),
			},
		},
	})
}

func getAgentAccessKey(limit int, reservation int, enabled bool) string {
	return fmt.Sprintf(`
resource "sysdig_agent_access_key" "my_agent_access_key" {
  limit       = %d
  reservation = %d
  enabled	  = %t
  metadata = {
    "test"             = "yes"
    "acceptance_test"  = "true"
  }
}

data "sysdig_agent_access_key" "data" {
  id = sysdig_agent_access_key.my_agent_access_key.id
}
`, limit, reservation, enabled)
}
