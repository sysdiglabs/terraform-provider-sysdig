//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_onprem_monitor || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccDataUser(t *testing.T) {
	randomSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getUser(randomSuffix),
			},
		},
	})
}

func getUser(suffix string) string {
	return fmt.Sprintf(`
resource "sysdig_user" "sample" {
  email = "terraform-test+user-%s@sysdig.com"
}

data "sysdig_user" "me" {
	depends_on = ["sysdig_user.sample"]
	email = sysdig_user.sample.email
}
`, suffix)
}
