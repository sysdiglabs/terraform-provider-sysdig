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

func TestAccRuleFilesystem(t *testing.T) {
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
				Config: ruleFilesystemWithName(rText()),
			},
			{
				Config: ruleFilesystemWithoutTagsWithName(rText()),
			},
			{
				Config: ruleFilesystemWithReadonlyWithName(rText()),
			},
			{
				Config: ruleFilesystemWithReadwriteWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_rule_filesystem.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ruleFilesystemMinimalConfig(rText()),
			},
		},
	})
}

func ruleFilesystemWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_filesystem"  "foo" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  tags = ["filesystem", "cis"]

  read_only {
    matching = true // default
    paths = ["/etc"]
  }

  read_write {
    matching = false // default
    paths = ["/tmp"]
  }
}`, name, name)
}

func ruleFilesystemWithoutTagsWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_filesystem"  "foo" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  
  read_only {
    matching = true // default
    paths = ["/etc"]
  }

  read_write {
    matching = false // default
    paths = ["/tmp"]
  }
}`, name, name)
}

func ruleFilesystemWithReadonlyWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_filesystem"  "foo" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  
  read_only {
    matching = true // default
    paths = ["/etc"]
  }
}`, name, name)
}

func ruleFilesystemWithReadwriteWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_filesystem"  "foo" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  
  read_write {
    matching = true // default
    paths = ["/etc"]
  }
}`, name, name)
}

func ruleFilesystemMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_filesystem"  "foo-minimal" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
}`, name, name)
}
