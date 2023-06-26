//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

func TestAccPosturePolicy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigIBMSecureAPIKeyEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: listAllPosturePolicies(),
			},
		},
	})
}

func listAllPosturePolicies() string {
	return `
data sysdig_secure_posture_policies policies {}
`
}
