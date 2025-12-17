//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

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

func TestAccRuleSyscall(t *testing.T) {
	t.Skip("List matching rules are deprecated - skipping tests")

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
				Config: ruleSyscallWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_rule_syscall.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: ruleSyscallWithoutTags(rText()),
			},
			{
				Config: ruleSyscallWithMinimalConfig(rText()),
			},
		},
	})
}

func ruleSyscallWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_syscall" "foo" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  tags = ["syscall", "cis"]

  matching = true // default
  syscalls = ["open", "execve"]
}`, name, name)
}

func ruleSyscallWithoutTags(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_syscall" "foo" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"

  matching = true // default
  syscalls = ["open", "execve"]
}`, name, name)
}

func ruleSyscallWithMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_syscall" "foo-minimal" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"

}`, name, name)
}
