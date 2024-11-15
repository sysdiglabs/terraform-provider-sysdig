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

func TestAccRuleFalcoCountDataSource(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

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
				Config: ruleFalcoCountDataSource(rText()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_rule_falco_count.data_terminal_shell", "rule_count", "1"),
				),
			},
			{
				Config: ruleFalcoCountDataSourceWithAppends(rText()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sysdig_secure_rule_falco_count.data_terminal_shell", "rule_count", "2"),
				),
			},
		},
	})
}

func ruleFalcoCountDataSource(name string) string {
	return fmt.Sprintf(`
%s

data "sysdig_secure_rule_falco_count" "data_terminal_shell" {
	name = "TERRAFORM TEST %s - Terminal Shell"
	depends_on = [ sysdig_secure_rule_falco.terminal_shell ]
}
`, ruleFalcoTerminalShell(name), name)
}

func ruleFalcoCountDataSourceWithAppends(name string) string {
	return fmt.Sprintf(`
%s

resource "sysdig_secure_rule_falco" "terminal_shell_append" {
	name = "TERRAFORM TEST %s - Terminal Shell"
  
	condition = "and never_true"
	source = "syscall" // syscall or k8s_audit
	append = true
	depends_on = [ sysdig_secure_rule_falco.terminal_shell ]
}

data "sysdig_secure_rule_falco_count" "data_terminal_shell" {
	name = "TERRAFORM TEST %s - Terminal Shell"
	depends_on = [ sysdig_secure_rule_falco.terminal_shell, sysdig_secure_rule_falco.terminal_shell_append ]
}
`, ruleFalcoTerminalShell(name), name, name)
}
