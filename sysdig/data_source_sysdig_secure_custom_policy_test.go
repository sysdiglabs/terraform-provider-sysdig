//go:build tf_acc_sysdig_secure || tf_acc_policies

package sysdig_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccCustomPolicyDataSource(t *testing.T) {
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
				Config: customPolicyDataSource(rText),
			},
		},
	})
}

func customPolicyDataSource(name string) string {
	return `
resource "sysdig_secure_custom_policy" "sample" {
	name = "%s"
	description = "Test Description"
	enabled = true
}
	
data "sysdig_secure_custom_policy" "example" {
	name = "%s"
	depends_on=[ sysdig_secure_custom_policy.sample ]
}
`
}
