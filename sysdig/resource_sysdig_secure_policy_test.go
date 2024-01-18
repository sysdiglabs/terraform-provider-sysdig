//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccPolicy(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: policyWithName(rText()),
			},
			{
				ResourceName:      "sysdig_secure_policy.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: policyWithoutActions(rText()),
			},
			{
				Config: policyWithoutNotificationChannels(rText()),
			},
			{
				Config: policyWithMinimumConfiguration(rText()),
			},
			{
				Config: policiesWithDifferentSeverities(rText()),
			},
			{
				Config: policiesWithKillAction(rText()),
			},
			{
				Config: policiesForAWSCloudtrail(rText()),
			},
			{
				Config: policiesForGCPAuditLog(rText()),
			},
			{
				Config: policiesForAzurePlatformlogs(rText()),
			},
		},
	})
}

func policyWithName(name string) string {
	return fmt.Sprintf(`
%s
%s
resource "sysdig_secure_policy" "sample" {
  name = "TERRAFORM TEST 1 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""
  rule_names = [sysdig_secure_rule_falco.terminal_shell.name]
  runbook = "https://sysdig.com"

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
      name = "testcapture"
      filter = "proc.name=cat"
      bucket_name = "testbucket"
      folder = "testfolder"
    }
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}
`, secureNotificationChannelEmailWithName(name), ruleFalcoTerminalShell(name), name, name)
}

func policyWithoutActions(name string) string {
	return fmt.Sprintf(`
%s
%s
resource "sysdig_secure_policy" "sample2" {
  name = "TERRAFORM TEST 2 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""
  rule_names = [sysdig_secure_rule_falco.terminal_shell.name]

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]

  actions {}
}
`, secureNotificationChannelEmailWithName(name), ruleFalcoTerminalShell(name), name, name)
}

func policyWithoutNotificationChannels(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_policy" "sample3" {
  name = "TERRAFORM TEST 3 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""
  rule_names = ["Terminal shell in container"]
  actions {}
}
`, name, name)
}

func policyWithMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_policy" "sample4" {
  name = "TERRAFORM TEST 4 %s"
  description = "TERRAFORM TEST %s"
  actions {}
}
`, name, name)
}

func policiesWithDifferentSeverities(name string) (res string) {
	for i := 0; i <= 7; i++ {
		res += fmt.Sprintf(`
resource "sysdig_secure_policy" "sample_%d" {
  name = "TERRAFORM TEST 1 %s-%d"
  description = "TERRAFORM TEST %s-%d"
  enabled = true
  severity = %d
  scope = "container.id != \"\""
  rule_names = ["Terminal shell in container"]

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
      name = "capture_name"
      filter = "proc.name=cat"
      bucket_name = "testbucket"
    }
  }
}

`, i, name, i, name, i, i)
	}
	return
}

func policiesWithKillAction(name string) (res string) {
	return fmt.Sprintf(`
resource "sysdig_secure_policy" "sample" {
  name = "TERRAFORM TEST 1 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""
  rule_names = ["Terminal shell in container"]

  actions {
    container = "kill"
  }
}
`, name, name)
}

func policiesForAWSCloudtrail(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_policy" "sample4" {
  name = "TERRAFORM TEST 4 %s"
  description = "TERRAFORM TEST %s"
  type = "aws_cloudtrail"
  actions {}
}
`, name, name)
}

func policiesForGCPAuditLog(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_policy" "sample5" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  type = "gcp_auditlog"
  actions {}
}
`, name, name)
}

func policiesForAzurePlatformlogs(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_policy" "sample6" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  type = "azure_platformlogs"
  actions {}
}
`, name, name)
}
