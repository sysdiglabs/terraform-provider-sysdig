//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

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

func TestAccMLPolicyDataSource(t *testing.T) {
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
				Config: mlPolicyDataSource(rText),
			},
		},
	})
}

func mlPolicyDataSource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_ml_policy" "policy_1" {
  name        = "Test ML Policy %s"
  description = "Test ML Policy Description %s"
  enabled     = true
  severity    = 4

  rule {
    description = "Test ML Rule Description"

    cryptomining_trigger {
      enabled   = true
      threshold = 1
      severity  = 1
    }
  }

}
	
data "sysdig_secure_ml_policy" "policy_2" {
  name       = sysdig_secure_ml_policy.policy_1.name
  depends_on = [sysdig_secure_ml_policy.policy_1]
}
`, name, name)
}
