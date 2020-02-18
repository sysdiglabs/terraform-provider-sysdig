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

func TestAccRuleSyscall(t *testing.T) {
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
				Config: ruleSyscallWithName(rText()),
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
