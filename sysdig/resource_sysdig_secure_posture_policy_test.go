//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSecurePosturePolicy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv, SysdigIBMSecureAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: createPolicyResource("test"),
			},
			{
				Config: updatePolicyResource("test"),
			},
		},
	})
}

func createPolicyResource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_policy" "sample" {
  name = "policy-test"
  description = "policy description"
}`, name)
}

func updatePolicyResource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_posture_policy" "sample" {
		name = "policy-test"
		description = "updated policy description"
}`, name)
}
