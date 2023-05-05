//go:build tf_acc_sysdig || tf_acc_sysdig_secure || tf_acc_policies

package sysdig_test

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
				Config: managedPoliciesWithKillAction(),
			},
		},
	})
}

func managedPolicy() {
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

func managedPolicyWithoutActions() {
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

func managedPolicyWithoutNotificationChannels() {
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

func managedPolicyWithMinimumConfiguration() {
	return fmt.Sprintf(`
resource "sysdig_secure_managed_policy" "sample" {
	name = "Sysdig Runtime Threat Detection"
	enabled = true
}
	`)
}

func managedPolicyWithKillAction() {
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
