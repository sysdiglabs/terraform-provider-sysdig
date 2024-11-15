//go:build tf_acc_sysdig || tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

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

func TestAccRuleFalcoDataSource(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	rTextForAppendTest := rText()

	resource.Test(t, resource.TestCase{
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
				Config: ruleFalcoDataSource(rText()),
			},
			{
				Config: setupRuleFalcoDataSourceWithAppends(rTextForAppendTest),
			},
			{
				Config: ruleFalcoDataSourceWithAppends(rTextForAppendTest),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sysdig_secure_rule_falco.data_terminal_shell.0", "id"),
					resource.TestCheckResourceAttrSet("data.sysdig_secure_rule_falco.data_terminal_shell.1", "id"),
				),
			},
		},
	})
}

func ruleFalcoDataSource(name string) string {
	return fmt.Sprintf(`
%s

data "sysdig_secure_rule_falco" "data_terminal_shell" {
	name = "TERRAFORM TEST %s - Terminal Shell"
	depends_on = [ sysdig_secure_rule_falco.terminal_shell ]
}
`, ruleFalcoTerminalShell(name), name)
}

func setupRuleFalcoDataSourceWithAppends(name string) string {
	return fmt.Sprintf(`
	%s

	resource "sysdig_secure_rule_falco" "terminal_shell_append" {
		name = "TERRAFORM TEST %s - Terminal Shell"

		condition = "and never_true"
		source = "syscall" // syscall or k8s_audit
		append = true
		depends_on = [ sysdig_secure_rule_falco.terminal_shell ]
	}
`, ruleFalcoTerminalShell(name), name)
}

func ruleFalcoDataSourceWithAppends(name string) string {
	return fmt.Sprintf(`
data "sysdig_secure_rule_falco_count" "terminal_shell_count" {
	name = "TERRAFORM TEST %s - Terminal Shell"
}

data "sysdig_secure_rule_falco" "data_terminal_shell" {
	count = data.sysdig_secure_rule_falco_count.terminal_shell_count.rule_count
	name = "TERRAFORM TEST %s - Terminal Shell"
	index = "${count.index}"

	depends_on = [ data.sysdig_secure_rule_falco_count.terminal_shell_count ]
}
`, name, name)
}
