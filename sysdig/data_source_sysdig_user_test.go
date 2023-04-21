//go:build tf_acc_sysdig

package sysdig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccDataUser(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv, SysdigIBMMonitorAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getUser(),
			},
		},
	})
}

func getUser() string {
	return `
resource "sysdig_user" "sample" {
  email = "terraform-test+user@sysdig.com"
}

data "sysdig_user" "me" {
	depends_on = ["sysdig_user.sample"]
	email = sysdig_user.sample.email
}
`
}
