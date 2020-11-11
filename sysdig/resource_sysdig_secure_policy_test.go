package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func TestAccPolicy(t *testing.T) {
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
				Config: policyWithName(rText()),
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

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
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
}
`, name, name)
}

func policyWithMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_policy" "sample4" {
  name = "TERRAFORM TEST 4 %s"
  description = "TERRAFORM TEST %s"
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
