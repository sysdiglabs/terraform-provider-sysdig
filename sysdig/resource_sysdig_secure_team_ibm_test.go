//go:build tf_acc_ibm

package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func TestAccSecureIBMTeam(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_IBM_SECURE_API_KEY"); v == "" {
				t.Fatal("SYSDIG_IBM_SECURE_API_KEY must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureIBMTeamWithName(rText()),
			},
			{
				Config: secureIBMTeamMinimumConfiguration(rText()),
			},
			{
				ResourceName:      "sysdig_secure_team.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureIBMTeamWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name               = "sample-%s"
  description        = "%s"
  scope_by           = "container"
  filter             = "container.image.repo = \"sysdig/agent\""
}
`, name, name)
}

func secureIBMTeamMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_team" "sample" {
  name      = "sample-%s"
}`, name)
}
