//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
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
				Config: createPolicyResource(),
			},
			{
				Config: createPolicyWithRiskResource(),
			},
		},
	})
}
func createPolicyResource() string {
	return `resource "sysdig_secure_posture_policy" "sample" {
		name = "policy-test"
		description = "policy description"
		is_active = true
		type = "kubernetes"
	}`
}
func createPolicyWithRiskResource() string {
	return `
resource "sysdig_secure_posture_policy" "sample" {
	name = "policy-test-with-group"
	description = "updated policy description"
	is_active = true
	type = "kubernetes"
	group {
  		name = "group 1"
		description = "group 1"
  		requirement {
  			name = "requirement 1"
	  		description = "requirement 1"
  			control {
				name = "Create Pods"
			}
		}
	}
}`
}
