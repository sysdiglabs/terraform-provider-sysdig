//go:build tf_acc_sysdig_secure || tf_acc_ibm_secure

package sysdig_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPosturePoliciesDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigIBMSecureAPIKeyEnv, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: "data sysdig_secure_posture_policies policies {}",
				Check: func(state *terraform.State) error {
					policiesRef := "data.sysdig_secure_posture_policies.policies"
					s, ok := state.RootModule().Resources[policiesRef]
					if !ok {
						return fmt.Errorf("%s not found", policiesRef)
					}
					numOfPolicies, err := strconv.Atoi(s.Primary.Attributes["policies.#"])
					if err != nil {
						return err
					}

					if numOfPolicies == 0 {
						return fmt.Errorf("missing policies")
					}
					return nil
				},
			},
		},
	})
}
