//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_policies_okta

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccOktaMLPolicyDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: oktaOktaMLPolicyDataSource(rText),
			},
		},
	})
}

func oktaOktaMLPolicyDataSource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_okta_ml_policy" "policy_1" {
  name        = "Test Okta ML Policy %s"
  description = "Test Okta ML Policy Description %s"
  enabled     = true
  severity    = 4

  rule {
    description = "Test Okta ML Rule Description"

    anomalous_console_login {
      enabled   = true
      threshold = 2
    }
  }

}

data "sysdig_secure_okta_ml_policy" "policy_2" {
  name       = sysdig_secure_okta_ml_policy.policy_1.name
  depends_on = [sysdig_secure_okta_ml_policy.policy_1]
}
`, name, name)
}
