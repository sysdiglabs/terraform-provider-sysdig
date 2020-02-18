package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccRuleProcess(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: ruleProcessWithName(rText()),
			},
			{
				Config: ruleProcessWithoutTags(rText()),
			},
			{
				Config: ruleProcessWithMinimalConfig(rText()),
			},
		},
	})
}

func ruleProcessWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_process" "foo" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  tags = ["container", "cis"]

  matching = true // default
  processes = ["bash"]
}`, name, name)
}

func ruleProcessWithoutTags(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_process" "foo" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"

  matching = true // default
  processes = ["bash"]
}`, name, name)
}

func ruleProcessWithMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_process" "foo-minimal" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
}`, name, name)
}
