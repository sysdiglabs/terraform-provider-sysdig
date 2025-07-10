//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

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

//func TestAccRuleFalco(t *testing.T) {
//	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
//
//	ruleRandomImmutableText := rText()
//
//	randomText := rText()
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
//				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
//			}
//		},
//		ProviderFactories: map[string]func() (*schema.Provider, error){
//			"sysdig": func() (*schema.Provider, error) {
//				return sysdig.Provider(), nil
//			},
//		},
//		Steps: []resource.TestStep{
//			//{
//	Config: ruleFalcoTerminalShell(ruleRandomImmutableText),
//},
//{
//	Config: ruleFalcoUpdatedTerminalShell(ruleRandomImmutableText),
//},
//{
//	Config: ruleFalcoTerminalShellWithMinimumEngineVersion(rText()),
//},
//{
//	ResourceName:      "sysdig_secure_rule_falco.terminal_shell",
//	ImportState:       true,
//	ImportStateVerify: true,
//},
//{
//	Config: ruleFalcoTerminalShellWithAppend(),
//},
//{
//	ResourceName:      "sysdig_secure_rule_falco.terminal_shell_append",
//	ImportState:       true,
//	ImportStateVerify: true,
//},
//{
//	Config: ruleFalcoGcpAuditlog(rText()),
//},
//{
//	Config: ruleFalcoAzureAuditlog(rText()),
//},
//{
//	Config: ruleFalcoKubeAudit(rText()),
//},
//{
//	ResourceName:      "sysdig_secure_rule_falco.kube_audit",
//	ImportState:       true,
//	ImportStateVerify: true,
//},
// Incorrect configurations
//{
//	Config:      ruleFalcoTerminalShellWithMissingOuput(rText()),
//	ExpectError: regexp.MustCompile("output must be set when append = false"),
//},
//{
//	Config:      ruleFalcoTerminalShellWithMissingSource(rText()),
//	ExpectError: regexp.MustCompile("source must be set when append = false"),
//},
//{
//	Config: ruleFalcoWithExceptions(randomText),
//},
//{
//	ResourceName:      "sysdig_secure_rule_falco.falco_rule_with_exceptions",
//	ImportState:       true,
//	ImportStateVerify: true,
//},
//{
//	Config: existingFalcoRuleWithExceptions(randomText),
//},
//{
//	ResourceName:      "sysdig_secure_rule_falco.attach_to_cluster_admin_role_exceptions",
//	ImportState:       true,
//	ImportStateVerify: true,
//},
//{
//	Config: ruleFalcoCloudAWSCloudtrail(randomText),
//},
//{
//	Config: ruleFalcoCloudAWSCloudtrailWithAppend(),
//},
//{
//	Config: ruleOkta(randomText),
//},
//{
//	Config: ruleOktaWithAppend(),
//},
//{
//	Config: ruleGithub(randomText),
//},
//{
//	Config: ruleGithubWithAppend(),
//},
//		},
//	})
//}

func TestAccRuleFalcoTerminalShell(t *testing.T) {
	ruleRandomImmutableText := randomString()
	steps := []resource.TestStep{
		{
			Config: ruleFalcoTerminalShell(ruleRandomImmutableText),
		},
		{
			Config: ruleFalcoUpdatedTerminalShell(ruleRandomImmutableText),
		},
		{
			ResourceName:      "sysdig_secure_rule_falco.terminal_shell",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}
	runTest(steps, t)
}

func TestAccRuleFalcoTerminalShellWithMinimumEngineVersion(t *testing.T) {
	steps := []resource.TestStep{
		{Config: ruleFalcoTerminalShellWithMinimumEngineVersion(randomString())},
	}
	runTest(steps, t)
}

func TestRuleFalcoTerminalShellWithAppend(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleFalcoTerminalShellWithAppend(),
		},
		{
			ResourceName:      "sysdig_secure_rule_falco.terminal_shell_append",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}
	runTest(steps, t)
}

func TestRuleFalcoGcpAuditlog(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleFalcoGcpAuditlog(randomString()),
		},
	}
	runTest(steps, t)
}

func TestRuleFalcoAzureAuditlog(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleFalcoAzureAuditlog(randomString()),
		},
	}
	runTest(steps, t)
}

func TestRuleFalcoKubeAudit(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleFalcoKubeAudit(randomString()),
		},
		{
			ResourceName:      "sysdig_secure_rule_falco.kube_audit",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}
	runTest(steps, t)
}

func TestIncorrectErrors(t *testing.T) {
	steps := []resource.TestStep{
		// Incorrect configurations
		{
			Config:      ruleFalcoTerminalShellWithMissingOuput(randomString()),
			ExpectError: regexp.MustCompile("output must be set when append = false"),
		},
		{
			Config:      ruleFalcoTerminalShellWithMissingSource(randomString()),
			ExpectError: regexp.MustCompile("source must be set when append = false"),
		},
	}
	runTest(steps, t)
}

func TestRuleFalcoWithExceptions(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleFalcoWithExceptions(randomString()),
		},
		{
			ResourceName:      "sysdig_secure_rule_falco.falco_rule_with_exceptions",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}
	runTest(steps, t)
}

func TestExistingFalcoRuleWithExceptions(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: existingFalcoRuleWithExceptions(randomString()),
		},
		{
			ResourceName:      "sysdig_secure_rule_falco.attach_to_cluster_admin_role_exceptions",
			ImportState:       true,
			ImportStateVerify: true,
		},
	}
	runTest(steps, t)
}

func TestRuleFalcoCloudAWSCloudtrail(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleFalcoCloudAWSCloudtrail(randomString()),
		},
	}
	runTest(steps, t)
}

func TestRuleFalcoCloudAWSCloudtrailAppend(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: RuleFalcoCloudAWSCloudtrailWithAppend(randomString()),
		},
	}
	runTest(steps, t)
}

func TestRuleOkta(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleOkta(randomString()),
		},
	}
	runTest(steps, t)
}

func TestRuleOktaAppends(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleOktaWithAppend(),
		},
	}
	runTest(steps, t)
}

func TestRuleGithub(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleGithub(randomString()),
		},
	}
	runTest(steps, t)
}

func TestRuleGithubAppends(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleGithubWithAppend(),
		},
	}
	runTest(steps, t)
}

func TestRuleGuardDuty(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleGuardDuty(randomString()),
		},
	}
	runTest(steps, t)
}

func TestRuleGuardDutyAppends(t *testing.T) {
	steps := []resource.TestStep{
		{
			Config: ruleGuardDutyWithAppend(),
		},
	}
	runTest(steps, t)
}

func randomString() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

func runTest(steps []resource.TestStep, t *testing.T) {
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
		Steps: steps,
	})
}

func ruleFalcoTerminalShell(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell" {
  name = "TERRAFORM TEST %s - Terminal Shell"
  description = "TERRAFORM TEST %s"
  tags = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output = "A shell was spawned in a container with an attached terminal (user=%%user.name %%container.info shell=%%proc.name parent=%%proc.pname cmdline=%%proc.cmdline terminal=%%proc.tty container_id=%%container.id image=%%container.image.repository)"
  priority = "notice"
  source = "syscall" // syscall or k8s_audit
}`, name, name)
}

func ruleFalcoTerminalShellWithMissingOuput(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell" {
  name = "TERRAFORM TEST %s - Terminal Shell Missing Output"
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
  name = "TERRAFORM TEST %s - Terminal Shell Missing Source"
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
  name = "TERRAFORM TEST %s - Terminal Shell Updated"
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

func ruleFalcoGcpAuditlog(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "gcp_audit" {
  name = "TERRAFORM TEST %[1]s - GCP Audit"
  description = "TERRAFORM TEST %[1]s"
  tags = ["gcp"]

  condition = "gcp.serviceName=\"compute.googleapis.com\" and gcp.methodName endswith \".compute.instances.setMetadata\""
  output = "GCP Audit Event received (%%gcp.serviceName, %%gcp.methodName)"
  priority = "debug"
  source = "gcp_auditlog"
}`, name)
}

func ruleFalcoAzureAuditlog(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "azure_audit" {
  name = "TERRAFORM TEST %[1]s - Azure Audit"
  description = "TERRAFORM TEST %[1]s"
  tags = ["azure"]

  condition = "jevt.value[/operationName] = \"DeleteBlob\""
  output = "Azure Audit Event received (%%jevt.value[/operationName])"
  priority = "debug"
  source = "azure_platformlogs"
}`, name)
}

func ruleFalcoTerminalShellWithAppend() string {
	return `
resource "sysdig_secure_rule_falco" "terminal_shell_append" {
  name = "Terminal shell in container" # Sysdig-provided
  condition = "and spawned_process and shell_procs and proc.tty != 0 and container_entrypoint"
  append = true
}`
}

func ruleFalcoWithExceptions(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "falco_rule_with_exceptions" {
  name        = "TERRAFORM TEST %s - Attach to cluster-admin Role"
  condition   = "kevt and clusterrolebinding and kcreate and ka.req.binding.role=cluster-admin"
  description = "Detect any attempt to create a ClusterRoleBinding to the cluster-admin user"

  output = "Cluster Role Binding to cluster-admin role (user=%%ka.user.name subject=%%ka.req.binding.subjects)"
  tags   = ["NIST_800-53_AC-2(12)(a)", "NIST_800-53_AU-6(8)", "NIST_800-53_SI-7(11)", "k8s", "SOC2_CC6.3", "NIST_800-53_AC-3", "NIST_800-53", "NIST_800-53_AC-2d", "SOC2"]
  source = "k8s_audit"

  exceptions {
   name = "subjects_with_in"
   fields = ["ka.req.binding.subjects", "ka.req.binding.role"]
   comps = ["in", "="]
   values = jsonencode([ [["sysdig", "sysdiglabs"], "falco"] ])
  }
  exceptions {
   name = "subjects_equal"
   fields = ["ka.req.binding.subjects", "ka.req.binding.role"]
   comps = ["=", "="]
   values = jsonencode([ ["foo", "bar"] ])
  }
  exceptions {
   name = "only_one_field"
   fields = ["ka.req.binding.subjects"]
   comps = ["in"]
   values = jsonencode([["foo"]])
  }
  exceptions {
   name = "only_one_field_without_comps"
   fields = ["ka.req.binding.subjects"]
   values = jsonencode([["foo"]])
  }
}
`, name)
}

func existingFalcoRuleWithExceptions(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "attach_to_cluster_admin_role_exceptions" {
    name = "Terminal shell in container" # Sysdig-provided
    append    = true

    exceptions {
        name = "proc_name_%s"
        fields = ["proc.name"]
        comps = ["in"]
        values = jsonencode([["sh"]])
   }
}`, name)
}

func ruleFalcoTerminalShellWithMinimumEngineVersion(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell" {
  name = "TERRAFORM TEST %s - Terminal Shell Min Engine Version"
  minimum_engine_version = 13
  description = "TERRAFORM TEST %s"
  tags = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output = "A shell was spawned in a container with an attached terminal (user=%%user.name %%container.info shell=%%proc.name parent=%%proc.pname cmdline=%%proc.cmdline terminal=%%proc.tty container_id=%%container.id image=%%container.image.repository)"
  priority = "notice"
  source = "syscall" // syscall or k8s_audit
}`, name, name)
}

func ruleFalcoCloudAWSCloudtrail(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "awscloudtrail" {
  name = "TERRAFORM TEST %[1]s - AWSCloudtrail"
  description = "TERRAFORM TEST %[1]s"
  tags = ["awscloudtrail"]

  condition = "ct.name=\"CreateApp\""
  output = "AWSCloudtrail Event received (requesting user=%%ct.user)"
  priority = "debug"
  source = "awscloudtrail"
}`, name, name)
}

func RuleFalcoCloudAWSCloudtrailWithAppend(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "awscloudtrail_append" {
  name = "Amplify Create App"
  source = "awscloudtrail"
  append = true
  exceptions {
	name = "user_name_%s"
	fields = ["ct.user"]
	comps = ["="]
	values = jsonencode([ ["user_a"] ])
   }
}`, name)
}

func ruleOkta(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "okta" {
  name = "TERRAFORM TEST %[1]s - Okta"
  description = "TERRAFORM TEST %[1]s"
  tags = ["okta"]

  condition = "okta.evt.type=\"user.account.update_password\""
  output = "Okta Event received (okta.severity=%%okta.severity)"
  priority = "debug"
  source = "okta"
}`, name, name)
}

func ruleOktaWithAppend() string {
	return `
resource "sysdig_secure_rule_falco" "okta_append" {
  name = "User changing password in to Okta"
  source = "okta"
  append = true
  exceptions {
	name = "actor_name"
	fields = ["okta.actor.name"]
	comps = ["="]
	values = jsonencode([ ["user_b"] ])
   }
}`
}

func ruleGithub(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "github" {
  name = "TERRAFORM TEST %[1]s - Github"
  description = "TERRAFORM TEST %[1]s"
  tags = ["github"]

  condition = "github.action=\"delete\""
  output = "Github Event received (github.user=%%github.user)"
  priority = "debug"
  source = "github"
}`, name, name)
}

func ruleGithubWithAppend() string {
	return `
resource "sysdig_secure_rule_falco" "github_append" {
  name = "Github Webhook Connected"
  source = "github"
  append = true
  exceptions {
	name = "user_name"
	fields = ["github.user"]
	comps = ["="]
	values = jsonencode([ ["user_c"] ])
   }
}`
}

func ruleGuardDuty(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_secure_rule_falco" "guardduty" {
	  name = "TERRAFORM TEST %[1]s - GuardDuty"
	  description = "TERRAFORM TEST %[1]s"
	  tags = ["guardduty"]
	
	  condition = "guardduty.resourceType=\"Container\""
	  output = "GuardDuty Event received (account ID=%%guardduty.accountId)"
	  priority = "debug"
	  source = "guardduty"
	}`, name, name)
}

func ruleGuardDutyWithAppend() string {
	return `
	resource "sysdig_secure_rule_falco" "guardduty_append" {
	  name = "GuardDuty High Severity Finding on Container"
	  source = "guardduty"
	  append = true
	  exceptions {
		name = "resource_type_tf"
		fields = ["guardduty.resourceType"]
		comps = ["="]
		values = jsonencode([ ["Amazon S2"] ])
	   }
	}`
}
