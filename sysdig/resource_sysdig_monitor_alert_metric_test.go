//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor || tf_acc_onprem_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/buildinfo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccAlertMetric(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigMonitorApiTokenEnv, SysdigIBMMonitorAPIKeyEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: alertMetricWithName(rText()),
			},
			{
				Config: alertMetricWithGroupName(rText()),
			},
			{
				Config: alertMetricWithoutScopeWithName(rText()),
			},
			{
				Config: alertMetricWithNotificationChannel(rText()),
				SkipFunc: func() (bool, error) {
					return buildinfo.IBMMonitor, nil
				},
			},
			{
				ResourceName:      "sysdig_secure_notification_channel_pagerduty.sample-pagerduty",
				ImportState:       true,
				ImportStateVerify: true,
				SkipFunc: func() (bool, error) {
					return buildinfo.IBMMonitor, nil
				},
			},
		},
	})
}

func alertMetricWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_metric" "sample" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 3

	metric = "avg(avg(cpu.used.percent)) > 50"
	scope = "agent.id in (\"foo\")"
	
	trigger_after_minutes = 10

	enabled = false

	multiple_alerts_by = ["kubernetes.deployment.name"]

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}

func alertMetricWithGroupName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_metric" "sample" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 3

	metric = "avg(avg(cpu.used.percent)) > 50"
	scope = "agent.id in (\"foo\")"
	
	trigger_after_minutes = 10
	group_name = "sample_group_name"
	enabled = false

	multiple_alerts_by = ["kubernetes.deployment.name"]

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}

func alertMetricWithoutScopeWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_metric" "sample2" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 3

	metric = "avg(avg(cpu.used.percent)) > 50"
	
	trigger_after_minutes = 10

	enabled = false

	multiple_alerts_by = ["host.hostName"]

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}

// Reported by @logdnalf at https://github.com/draios/terraform-provider-sysdig/issues/24
func alertMetricWithNotificationChannel(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_pagerduty" "sample-pagerduty" {
	name = "Example Channel %s - Pagerduty"
	enabled = true
	account = "account"
	service_key = "XXXXXXXXXX"
	service_name = "sysdig"
	notify_when_ok = true
	notify_when_resolved = true
}

resource "sysdig_monitor_alert_metric" "sample3" {
	enabled = true
	name = "TERAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 6
	metric = "sum(min(cpu.used.percent)) > 100000"
	scope = "agent.id in (\"foo\")"
	trigger_after_minutes = 20
	notification_channels = [
	sysdig_secure_notification_channel_pagerduty.sample-pagerduty.id
	]
	multiple_alerts_by = [
	"host.hostName"
	]
}`, name, name, name)
}
