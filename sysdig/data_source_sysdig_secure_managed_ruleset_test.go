//go:build tf_acc_sysdig || tf_acc_sysdig_secure

package sysdig_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccManagedPolicyDataSource(t *testing.T) {
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
				Config: managedRulesetDataSource(rText),
			},
		},
	})
}

func managedRulesetDataSource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_managed_ruleset" "sample" {
	name = "%s"
	description = "Test Description"
	inherited_from {
		name = "Sysdig Runtime Threat Detection"
		type = "falco"
	}
	enabled = true
}

data "sysdig_secure_managed_ruleset" "example" {
	depends_on = [sysdig_secure_managed_ruleset.sample]
	name = "%s"
	type = "falco
}
`, name, name)
}
