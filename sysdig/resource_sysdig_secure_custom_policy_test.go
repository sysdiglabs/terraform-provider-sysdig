//go:build tf_acc_sysdig_secure || tf_acc_policies || tf_acc_onprem_secure

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/buildinfo"
	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccCustomPolicy(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	policy1 := rText()

	steps := []resource.TestStep{
		{
			Config: customPolicyWithName(policy1),
		},
		{
			ResourceName:      "sysdig_secure_custom_policy.sample",
			ImportState:       true,
			ImportStateVerify: true,
		},
		{
			Config: customPolicyWithRulesOrderChange(policy1),
		},
		{
			Config: customPolicyWithoutActions(rText()),
		},
		{
			Config: customPolicyWithoutNotificationChannels(rText()),
		},
		{
			Config: customPolicyWithMinimumConfiguration(rText()),
		},
		{
			Config: customPoliciesWithDifferentSeverities(rText()),
		},
		{
			Config: customPoliciesWithKillAction(rText()),
		},
		{
			Config: customPoliciesWithDisabledRules(rText()),
		},
		{
			Config: customPoliciesWithKillProcessAction(rText()),
		},
	}

	if !buildinfo.OnpremSecure {
		steps = append(steps,
			resource.TestStep{Config: customPoliciesForAWSCloudtrail(rText())},
			resource.TestStep{Config: customPoliciesForGCPAuditLog(rText())},
			resource.TestStep{Config: customPoliciesForAzurePlatformlogs(rText())},
		)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: steps,
	})
}

func customPolicyWithName(name string) string {
	return fmt.Sprintf(`
%s
%s
resource "sysdig_secure_custom_policy" "sample" {
  name = "TERRAFORM TEST 1 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""
  runbook = "https://sysdig.com"

  rules {
    name = "Write below etc"
    enabled = true
  }
  rules {
    name = sysdig_secure_rule_falco.terminal_shell.name
    enabled = true
  }

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
      name = "testcapture"
    }
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}
`, secureNotificationChannelEmailWithName(name), ruleFalcoTerminalShell(name), name, name)
}

func customPolicyWithRulesOrderChange(name string) string {
	return fmt.Sprintf(`
%s
%s
resource "sysdig_secure_custom_policy" "sample" {
  name = "TERRAFORM TEST 1 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""
  runbook = "https://sysdig.com"

  rules {
    name = sysdig_secure_rule_falco.terminal_shell.name
    enabled = true
  }
  rules {
    name = "Write below etc"
    enabled = true
  }

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
      name = "testcapture"
    }
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}
`, secureNotificationChannelEmailWithName(name), ruleFalcoTerminalShell(name), name, name)
}

func customPolicyWithoutActions(name string) string {
	return fmt.Sprintf(`
%s
%s
resource "sysdig_secure_custom_policy" "sample2" {
  name = "TERRAFORM TEST 2 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]

  rules {
    name = sysdig_secure_rule_falco.terminal_shell.name
    enabled = true
  }

  actions {}
}
`, secureNotificationChannelEmailWithName(name), ruleFalcoTerminalShell(name), name, name)
}

func customPolicyWithoutNotificationChannels(name string) string {
	return fmt.Sprintf(`
%s
resource "sysdig_secure_custom_policy" "sample3" {
  name = "TERRAFORM TEST 3 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""

  rules {
    name = sysdig_secure_rule_falco.terminal_shell.name
    enabled = true
  }
  actions {}
}
`, ruleFalcoTerminalShell(name), name, name)
}

func customPolicyWithMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_custom_policy" "sample4" {
  name = "TERRAFORM TEST 4 %s"
  description = "TERRAFORM TEST %s"
  actions {}
}
`, name, name)
}

func customPoliciesWithDifferentSeverities(name string) (res string) {
	for i := 0; i <= 7; i++ {
		res += fmt.Sprintf(`
resource "sysdig_secure_custom_policy" "sample_%d" {
  name = "TERRAFORM TEST 1 %s-%d"
  description = "TERRAFORM TEST %s-%d"
  enabled = true
  severity = %d
  scope = "container.id != \"\""
  rules {
    name = "Terminal shell in container"
    enabled = true
  }

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
      name = "capture_name"
    }
  }
}

`, i, name, i, name, i, i)
	}
	return
}

func customPoliciesWithKillAction(name string) (res string) {
	return fmt.Sprintf(`
resource "sysdig_secure_custom_policy" "sample10" {
  name = "TERRAFORM TEST 10 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""

  rules {
    name = "Terminal shell in container"
    enabled = true
  }

  actions {
    container = "kill"
  }
}
`, name, name)
}

func customPoliciesWithKillProcessAction(name string) (res string) {
	return fmt.Sprintf(`
resource "sysdig_secure_custom_policy" "sample10" {
 name = "TERRAFORM TEST 1 %s"
 description = "TERRAFORM TEST %s"
 enabled = true
 severity = 4
 scope = "container.id != \"\""

 rules {
   name = "Terminal shell in container"
   enabled = true
 }

 actions {
   kill_process = "true"
 }
}
`, name, name)
}

func customPoliciesForAWSCloudtrail(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_custom_policy" "sample4" {
  name = "TERRAFORM TEST 4 %s"
  description = "TERRAFORM TEST %s"
  type = "aws_cloudtrail"
  actions {}
}
`, name, name)
}

func customPoliciesForGCPAuditLog(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_custom_policy" "sample5" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  type = "gcp_auditlog"
  actions {}
}
`, name, name)
}

func customPoliciesForAzurePlatformlogs(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_custom_policy" "sample6" {
  name = "TERRAFORM TEST %s"
  description = "TERRAFORM TEST %s"
  type = "azure_platformlogs"
  actions {}
}
`, name, name)
}

func customPoliciesWithDisabledRules(name string) string {
	return fmt.Sprintf(`
%s
%s
resource "sysdig_secure_custom_policy" "sample" {
  name = "TERRAFORM TEST 1 %s"
  description = "TERRAFORM TEST %s"
  enabled = true
  severity = 4
  scope = "container.id != \"\""
  runbook = "https://sysdig.com"

  rules {
    name = sysdig_secure_rule_falco.terminal_shell.name
    enabled = false
  }

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
      name = "testcapture"
    }
  }

  notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}
`, secureNotificationChannelEmailWithName(name), ruleFalcoTerminalShell(name), name, name)
}
