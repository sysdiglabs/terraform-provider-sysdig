//go:build tf_acc_sysdig_monitor || tf_acc_ibm_monitor

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccAlertV2FormBasedPrometheusTest(t *testing.T) {
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
				Config: alertV2FormBasedPrometheusTest(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithNoData(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithNotificationChannels(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithDescription(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithSeverity(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithGroup(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithCustomNotifications(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithLink(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithEnabled(rText()),
			},
			{
				Config: alertV2FormBasedPrometheusTestWithWarningThreshold(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_v2_form_based_prometheus.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertV2FormBasedPrometheusTest(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
}
`, name)
}

func alertV2FormBasedPrometheusTestWithNoData(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	no_data_behaviour = "TRIGGER"
}
`, name)
}

func alertV2FormBasedPrometheusTestWithNotificationChannels(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "nc_email_1" {
	name = "%s1"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_notification_channel_email" "nc_email_2" {
	name = "%s2"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {

	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	enabled = false
	notification_channels {
		id = sysdig_monitor_notification_channel_email.nc_email_1.id
		notify_on_resolve = false
	}
	notification_channels {
		id = sysdig_monitor_notification_channel_email.nc_email_2.id
		renotify_every_minutes = 30
	}
}
`, name, name, name)
}

func alertV2FormBasedPrometheusTestWithDescription(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	description = "description"
}
`, name)
}

func alertV2FormBasedPrometheusTestWithSeverity(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	severity = "high"
}
`, name)
}

func alertV2FormBasedPrometheusTestWithGroup(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	group = "customgroup"
}
`, name)
}

func alertV2FormBasedPrometheusTestWithCustomNotifications(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	custom_notification {
		subject = "test"
		prepend = "pre"
		append = "post"
	}
}
`, name)
}

func alertV2FormBasedPrometheusTestWithLink(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard_1" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	description = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart"
		name = "example panel"
		description = "description"

        legend {
            show_current = true
            position = "bottom"
            layout = "inline"
        }

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"

            format {
                display_format = "auto"
                input_format = "0-100"
                y_axis = "auto"
                null_value_display_mode = "nullGap"
            }
		}
	}
}

resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	link {
		type = "runbook"
		href = "http://example.com"
	}
	link {
		type = "dashboard"
		id = sysdig_monitor_dashboard.dashboard_1.id
	}
}
`, name, name, name)
}

func alertV2FormBasedPrometheusTestWithEnabled(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	enabled = false
}
`, name)
}

func alertV2FormBasedPrometheusTestWithWarningThreshold(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "nc_email_3" {
	name = "%s3"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_notification_channel_email" "nc_email_4" {
	name = "%s4"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {

	name = "TERRAFORM TEST - FORM BASED PROMETHEUS %s"
	query = "avg_over_time(sysdig_container_cpu_used_percent{container_name=\"test\"}[59s])"
	operator = ">="
	threshold = 50
	enabled = false
	warning_threshold = 10
	notification_channels {
		id = sysdig_monitor_notification_channel_email.nc_email_3.id
	}
	notification_channels {
		id = sysdig_monitor_notification_channel_email.nc_email_4.id
		warning_threshold = true
	}
}
`, name, name, name)
}
