//go:build sysdig_secure || tf_acc_policies

package sysdig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccRuleContainer(t *testing.T) {
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
				Config: ruleContainerWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_rule_container.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ruleContainerWithNameWithoutTags(rText()),
			},
			{
				Config: ruleContainerWithNameAndMinimumConfig(rText()),
			},
		},
	})
}

func ruleContainerWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_container" "sample" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  tags = ["container", "cis"]

  matching = true // default
  containers = ["foo", "foo:bar"]
}`, name, name)
}

func ruleContainerWithNameWithoutTags(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_container" "sample" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"

  matching = true // default
  containers = ["foo", "foo:bar"]
}`, name, name)
}

func ruleContainerWithNameAndMinimumConfig(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_container" "sample2" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
}`, name, name)
}
