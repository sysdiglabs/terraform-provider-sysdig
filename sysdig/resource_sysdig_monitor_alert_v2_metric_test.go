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

func TestAccAlertV2Metric(t *testing.T) {
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
				Config: alertV2Metric(rText()),
			},
			{
				Config: alertV2MetricWithScope(rText()),
			},
			{
				Config: alertV2MetricWithNoData(rText()),
			},
			{
				Config: alertV2MetricWithNotificationChannels(rText()),
			},
			{
				Config: alertV2MetricWithDescription(rText()),
			},
			{
				Config: alertV2MetricWithSeverity(rText()),
			},
			{
				Config: alertV2MetricWithGroupBy(rText()),
			},
			{
				Config: alertV2MetricWithGroup(rText()),
			},
			{
				Config: alertV2MetricWithCustomNotifications(rText()),
			},
			{
				Config: alertV2MetricWithCapture(rText()),
			},
			{
				Config: alertV2MetricWithLink(rText()),
			},
			{
				Config: alertV2MetricWithEnabled(rText()),
			},
			{
				Config: alertV2MetricWithWarningThreshold(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_v2_metric.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertV2Metric(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15

}
`, name)
}

func alertV2MetricWithScope(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_alert_v2_metric" "sample" {

		name = "TERRAFORM TEST - METRICV2 %s"
		metric = "sysdig_container_cpu_used_percent"
		group_aggregation = "avg"
		time_aggregation = "avg"
		operator = ">="
		threshold = 50
		trigger_after_minutes = 15
		scope {
			label = "kube_cluster_name"
			operator = "in"
			values = ["thom-cluster1", "demo-env-prom"]
		}
		scope {
			label = "kube_cluster_name"
			operator = "equals"
			values = ["thom-cluster3"]
		}

	}
	`, name)
}

func alertV2MetricWithNoData(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	no_data_behaviour = "TRIGGER"

}
`, name)
}

func alertV2MetricWithNotificationChannels(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "nc_email1" {
	name = "%s1"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_notification_channel_email" "nc_email2" {
	name = "%s2"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	enabled = false
	notification_channels {
		id = sysdig_monitor_notification_channel_email.nc_email1.id
		notify_on_resolve = false
	}
	notification_channels {
		id = sysdig_monitor_notification_channel_email.nc_email2.id
		renotify_every_minutes = 30
	}
}
`, name, name, name)
}

func alertV2MetricWithDescription(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	description = "description"

}
`, name)
}

func alertV2MetricWithSeverity(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	severity = "high"

}
`, name)
}

func alertV2MetricWithGroupBy(name string) string {
	return fmt.Sprintf(`
	resource "sysdig_monitor_alert_v2_metric" "sample" {

		name = "TERRAFORM TEST - METRICV2 %s"
		metric = "sysdig_container_cpu_used_percent"
		group_aggregation = "avg"
		time_aggregation = "avg"
		operator = ">="
		threshold = 50
		trigger_after_minutes = 15
		group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]

	}
	`, name)
}

func alertV2MetricWithGroup(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	group = "customgroup"

}
`, name)
}

func alertV2MetricWithCustomNotifications(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	custom_notification {
		subject = "test"
		prepend = "pre"
		append = "post"
	}

}
`, name)
}

func alertV2MetricWithCapture(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	capture {
		filename = "test.scap"
	}
}
`, name)
}

func alertV2MetricWithLink(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"

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

resource "sysdig_monitor_alert_v2_metric" "sample" {
	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	link {
		type = "runbook"
		href = "http://ciao2.com"
	}
	link {
		type = "dashboard"
		id = sysdig_monitor_dashboard.dashboard.id
	}
}
`, name, name, name)
}

func alertV2MetricWithEnabled(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	enabled = false

}
`, name)
}

func alertV2MetricWithWarningThreshold(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "nc_email1" {
	name = "%s1"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_notification_channel_email" "nc_email2" {
	name = "%s2"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_alert_v2_metric" "sample" {

	name = "TERRAFORM TEST - METRICV2 %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	trigger_after_minutes = 15
	enabled = false
	warning_threshold = 10
	notification_channels {
		id = sysdig_monitor_notification_channel_email.nc_email1.id
	}
	notification_channels {
		id = sysdig_monitor_notification_channel_email.nc_email2.id
		warning_threshold = true
	}
}
`, name, name, name)
}
