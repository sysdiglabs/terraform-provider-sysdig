//go:build tf_acc_sysdig || tf_acc_sysdig_secure || tf_acc_policies

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
				Config: managedPolicy(),
			},
			{
				ResourceName:      "sysdig_secure_managed_policy.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: managedPolicyWithoutActions(),
			},
			{
				Config: managedPolicyWithoutNotificationChannels(),
			},
			{
				Config: managedPolicyWithMinimumConfiguration(),
			},
			{
				Config: managedPolicyWithKillAction(),
			},
		},
	})
}

func managedPolicy() string {
	return fmt.Sprintf(`
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
	scope = "container.id != \"\""
	disabled_rules = ["Suspicious Cron Modification"]
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
	`)
}

func managedPolicyWithoutActions() string {
	return fmt.Sprintf(`
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
	scope = "container.id != \"\""
	disabled_rules = ["Suspicious Cron Modification"]
	runbook = "https://sysdig.com"

	actions {}
	
	notification_channels = [sysdig_secure_notification_channel_email.sample_email.id]
}
	`)
}

func managedPolicyWithoutNotificationChannels() string {
	return fmt.Sprintf(`
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
	scope = "container.id != \"\""
	disabled_rules = ["Suspicious Cron Modification"]
	runbook = "https://sysdig.com"

	actions {
		container = "stop"
		capture {
		  seconds_before_event = 5
		  seconds_after_event = 10
		  name = "testcapture"
		}
	}	
}
	`)
}

func managedPolicyWithMinimumConfiguration() string {
	return fmt.Sprintf(`
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
}
	`)
}

func managedPolicyWithKillAction() string {
	return fmt.Sprintf(`
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
	scope = "container.id != \"\""
	disabled_rules = ["Suspicious Cron Modification"]
	runbook = "https://sysdig.com"

	actions {
		container = "kill"
	}
}
	`)
}
