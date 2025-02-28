//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccManagedPolicyDataSource(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: managedPolicyDataSource(),
		},
	}

	if !strings.HasSuffix(os.Getenv("SYSDIG_SECURE_URL"), "ibm.com") {
		steps = append(steps, resource.TestStep{
			Config: managedStatefulPolicyDataSource(),
		},
		)
	}
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
		Steps: steps,
	})
}

func managedPolicyDataSource() string {
	return `
data "sysdig_secure_managed_policy" "example" {
	name = "Sysdig Runtime Threat Detection"
	type = "falco"
}
`
}

func managedStatefulPolicyDataSource() string {
	return `
data "sysdig_secure_managed_policy" "stateful_example" {
	name = "Sysdig AWS Behavioral Analytics Threat Detection"
	enabled = false
	type = "awscloudtrail_stateful"
}
`
}
