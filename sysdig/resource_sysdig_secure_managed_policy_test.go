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

func TestAccManagedPolicy(t *testing.T) {
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
				Config: managedPolicyWithoutNotificationChannels(),
			},
			{
				Config: managedPolicyWithoutActions(rText()),
			},
			{
				Config: managedPolicy(rText()),
			},
			{
				Config: managedPolicyWithMinimumConfiguration(),
			},
			{
				Config:  managedPolicyWithKillAction(),
				Destroy: true,
			},
		},
	})
}

func managedPolicy(name string) string {
	return fmt.Sprintf(`
%s
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
	scope = "container.id != \"\""
	disabled_rules = ["Hexadecimal string detected"]
	runbook = "https://sysdig.com"

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
	`, secureNotificationChannelEmailWithName(name))
}

func managedPolicyWithoutActions(name string) string {
	return fmt.Sprintf(`
%s
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
	scope = "container.id != \"\""
	disabled_rules = ["Hexadecimal string detected"]
	runbook = "https://sysdig.com"

	actions {}
	
	notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}
	`, secureNotificationChannelEmailWithName(name))
}

func managedPolicyWithoutNotificationChannels() string {
	return `
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
	scope = "container.id != \"\""
	disabled_rules = ["Hexadecimal string detected"]
	runbook = "https://sysdig.com"

	actions {
		container = "stop"
		capture {
		  seconds_before_event = 5
		  seconds_after_event = 10
		  name = "testcapture"
		}
	}	
}`
}

func managedPolicyWithMinimumConfiguration() string {
	return `
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
}`
}

func managedPolicyWithKillAction() string {
	return `resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
	scope = "container.id != \"\""
	disabled_rules = ["Hexadecimal string detected"]
	runbook = "https://sysdig.com"

	actions {
		container = "kill"
	}
}`
}
