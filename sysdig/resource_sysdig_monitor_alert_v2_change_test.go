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

func TestAccAlertV2Change(t *testing.T) {
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
				Config: alertV2Change(rText()),
			},
			{
				Config: alertV2ChangeWithScope(rText()),
			},
			{
				Config: alertV2ChangeWithNotificationChannels(rText()),
			},
			{
				Config: alertV2ChangeWithDescription(rText()),
			},
			{
				Config: alertV2ChangeWithSeverity(rText()),
			},
			{
				Config: alertV2ChangeWithGroupBy(rText()),
			},
			{
				Config: alertV2ChangeWithGroup(rText()),
			},
			{
				Config: alertV2ChangeWithCustomNotifications(rText()),
			},
			{
				Config: alertV2ChangeWithLink(rText()),
			},
			{
				Config: alertV2ChangeWithEnabled(rText()),
			},
			{
				Config: alertV2ChangeWithWarningThreshold(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_v2_change.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertV2Change(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
}
`, name)
}

func alertV2ChangeWithScope(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
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

func alertV2ChangeWithNotificationChannels(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "nc_email1" {
	name = "%s1"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_notification_channel_email" "nc_email2" {
	name = "%s2"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
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

func alertV2ChangeWithDescription(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
	description = "description"
}
`, name)
}

func alertV2ChangeWithSeverity(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
	severity = "high"
}
`, name)
}

func alertV2ChangeWithGroupBy(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
}
	`, name)
}

func alertV2ChangeWithGroup(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
	group = "customgroup"
}
`, name)
}

func alertV2ChangeWithCustomNotifications(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
	custom_notification {
		subject = "test"
		prepend = "pre"
		append = "post"
	}
}
`, name)
}

func alertV2ChangeWithLink(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "TERRAFORM TEST - CHANGE %s"
	description = "TERRAFORM TEST - CHANGE %s"

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

resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
	link {
		type = "runbook"
		href = "http://example.com"
	}
	link {
		type = "dashboard"
		id = sysdig_monitor_dashboard.dashboard.id
	}
}
`, name, name, name)
}

func alertV2ChangeWithEnabled(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
	enabled = false
}
`, name)
}

func alertV2ChangeWithWarningThreshold(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "nc_email1" {
	name = "%s1"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_notification_channel_email" "nc_email2" {
	name = "%s2"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_alert_v2_change" "sample" {
	name = "TERRAFORM TEST - CHANGE %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	operator = ">="
	threshold = 50
	shorter_time_range_seconds = 300
	longer_time_range_seconds = 3600
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
