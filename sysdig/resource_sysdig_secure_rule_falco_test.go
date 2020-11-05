package sysdig_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccRuleFalco(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	ruleRandomImmutableText := rText()

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
				Config: ruleFalcoTerminalShell(ruleRandomImmutableText),
			},
			{
				Config: ruleFalcoUpdatedTerminalShell(ruleRandomImmutableText),
			},
			{
				Config: ruleFalcoTerminalShellWithAppend(),
			},
			{
				Config: ruleFalcoKubeAudit(rText()),
			},
			// Incorrect configurations
			{
				Config:      ruleFalcoTerminalShellWithMissingOuput(rText()),
				ExpectError: regexp.MustCompile("output must be set when append = false"),
			},
			{
				Config:      ruleFalcoTerminalShellWithMissingSource(rText()),
				ExpectError: regexp.MustCompile("source must be set when append = false"),
			},
		},
	})
}

func ruleFalcoTerminalShell(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell" {
  name = "TERRAFORM TEST %s - Terminal Shell"
  tags = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output = "A shell was spawned in a container with an attached terminal (user=%%user.name %%container.info shell=%%proc.name parent=%%proc.pname cmdline=%%proc.cmdline terminal=%%proc.tty container_id=%%container.id image=%%container.image.repository)"
  priority = "notice"
  source = "syscall" // syscall or k8s_audit
}`, name)
}

func ruleFalcoTerminalShellWithMissingOuput(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell" {
  name = "TERRAFORM TEST %s - Terminal Shell"
  description = "TERRAFORM TEST %s"
  tags = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  priority = "notice"
  source = "syscall" // syscall or k8s_audit
}`, name, name)
}

func ruleFalcoTerminalShellWithMissingSource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell" {
  name = "TERRAFORM TEST %s - Terminal Shell"
  description = "TERRAFORM TEST %s"
  tags = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output = "A shell was spawned in a container with an attached terminal (user=%%user.name %%container.info shell=%%proc.name parent=%%proc.pname cmdline=%%proc.cmdline terminal=%%proc.tty container_id=%%container.id image=%%container.image.repository)"
  priority = "notice"
  append = false
}`, name, name)
}

func ruleFalcoUpdatedTerminalShell(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell" {
  name = "TERRAFORM TEST %s - Terminal Shell"
  description = "TERRAFORM TEST %s"
  tags = ["shell", "mitre_execution"]

  condition = "spawned_process and shell_procs and proc.tty != 0 and container_entrypoint"
  output = "A shell was spawned in a container with an attached terminal (user=%%user.name %%container.info shell=%%proc.name parent=%%proc.pname cmdline=%%proc.cmdline terminal=%%proc.tty container_id=%%container.id image=%%container.image.repository)"
  priority = "notice"
  source = "syscall" // syscall or k8s_audit
}`, name, name)
}

func ruleFalcoKubeAudit(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "kube_audit" {
  name = "TERRAFORM TEST %s - KubeAudit"
  description = "TERRAFORM TEST %s"
  tags = ["k8s"]

  condition = "kall"
  output = "K8s Audit Event received (user=%%ka.user.name verb=%%ka.verb uri=%%ka.uri obj=%%jevt.obj)"
  priority = "debug"
  source = "k8s_audit" // syscall or k8s_audit
}`, name, name)
}

func ruleFalcoTerminalShellWithAppend() string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell_append" {
  name = "Terminal shell in container" # Sysdig-provided
  condition = "and spawned_process and shell_procs and proc.tty != 0 and container_entrypoint"
  append = true
}`)
}
