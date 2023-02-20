package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccScanningPolicy(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

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
				Config: scanningPolicyWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_scanning_policy.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func scanningPolicyWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_scanning_policy" "sample" {
  name = "TERRAFORM TEST 1 %s"
  comment = "TERRAFORM TEST %s"

  rules {
    gate = "dockerfile"
    trigger = "effective_user"
    action = "WARN"
    params {
        name = "users"
        value = "docker"
    }
    params {
        name = "type"
        value = "blacklist"
    }
  }
}
`, name, name)
}
