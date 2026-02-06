//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_ibm_monitor || tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccCurrentUser(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv, SysdigIBMMonitorAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getCurrentUser(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sysdig_current_user.me", "customer_id"),
					resource.TestCheckResourceAttrSet("data.sysdig_current_user.me", "customer_name"),
					resource.TestCheckResourceAttrSet("data.sysdig_current_user.me", "customer_external_id"),
				),
			},
		},
	})
}

func getCurrentUser() string {
	return `
data "sysdig_current_user" "me" {
}
`
}
