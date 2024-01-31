//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccAWSMLPolicyDataSource(t *testing.T) {
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
				Config: awsAWSMLPolicyDataSource(rText),
			},
		},
	})
}

func awsAWSMLPolicyDataSource(name string) string {
	return `
resource "sysdig_secure_aws_ml_policy" "policy_1" {
  name        = "Test AWS ML Policy 81z7b1xng6"
  description = "Test AWS ML Policy Description"
  enabled     = true
  severity    = 4

  rule {
    description = "Test AWS ML Rule Description"

    anomalous_console_login {
      enabled   = true
      threshold = 2
      severity  = 1
    }
  }

}
	
data "sysdig_secure_aws_ml_policy" "policy_2" {
  name       = sysdig_secure_aws_ml_policy.policy_1.name
  depends_on = [sysdig_secure_aws_ml_policy.policy_1]
}
`
}
