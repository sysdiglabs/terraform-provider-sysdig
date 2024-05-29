//go:build tf_acc_sysdig_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPosturePolicyDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: posturePolicyDataSource(),
				Check: func(state *terraform.State) error {
					policyRef := "data.sysdig_secure_posture_policy.policy"
					s, ok := state.RootModule().Resources[policyRef]
					if !ok {
						return fmt.Errorf("%s not found", policyRef)
					}
					if s.Primary.ID != "2" {
						return fmt.Errorf("expected policy ID to be 2, got %s", s.Primary.ID)
					}
					if s.Primary.Attributes["name"] != "Sysdig Kubernetes" {
						return fmt.Errorf("expected policy name to be `Sysdig Kubernetes`, got %s", s.Primary.Attributes["name"])
					}
					return nil
				},
			},
		},
	})
}

func posturePolicyDataSource() string {
	return fmt.Sprintf(`
data "sysdig_secure_posture_policy" "policy" {
	id = 2
}`)
}
