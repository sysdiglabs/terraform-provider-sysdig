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

func TestAccAlertV2GroupOutlier(t *testing.T) {
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
				Config: alertV2GroupOutlier(rText()),
			},
			{
				Config: alertV2GroupOutlierWithMAD(rText()),
			},
			{
				Config: alertV2GroupOutlierWithNoData(rText()),
			},
			{
				Config: alertV2GroupOutlierWithNotificationChannels(rText()),
			},
			{
				Config: alertV2GroupOutlierWithDescription(rText()),
			},
			{
				Config: alertV2GroupOutlierWithSeverity(rText()),
			},
			{
				Config: alertV2GroupOutlierWithGroup(rText()),
			},
			{
				Config: alertV2GroupOutlierWithCustomNotifications(rText()),
			},
			{
				Config: alertV2GroupOutlierWithCapture(rText()),
			},
			{
				Config: alertV2GroupOutlierWithLink(rText()),
			},
			{
				Config: alertV2GroupOutlierWithEnabled(rText()),
			},
			{
				Config: alertV2GroupOutlierWithUnreportedAlertNotificationsRetentionSec(rText()),
			},
			{
				ResourceName:      "sysdig_monitor_alert_v2_group_outlier.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func alertV2GroupOutlier(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1.2
	observation_window_minutes = 15

}
`, name)
}

func alertV2GroupOutlierWithMAD(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "MAD"
  mad_threshold = 10.1
  mad_tolerance = 5.5
	observation_window_minutes = 15

}
`, name)
}

func alertV2GroupOutlierWithNoData(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
	no_data_behaviour = "TRIGGER"

}
`, name)
}

func alertV2GroupOutlierWithNotificationChannels(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_notification_channel_email" "nc_email1" {
	name = "%s1"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_notification_channel_email" "nc_email2" {
	name = "%s2"
	recipients = ["root@localhost.com"]
}

resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
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

func alertV2GroupOutlierWithDescription(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
	description = "description"

}
`, name)
}

func alertV2GroupOutlierWithSeverity(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
	severity = "high"

}
`, name)
}

func alertV2GroupOutlierWithGroup(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
	group = "customgroup"

}
`, name)
}

func alertV2GroupOutlierWithCustomNotifications(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
	custom_notification {
		subject = "test"
		prepend = "pre"
		append = "post"
	}

}
`, name)
}

func alertV2GroupOutlierWithCapture(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
	capture {
		filename = "test.scap"
	}
}
`, name)
}

func alertV2GroupOutlierWithLink(name string) string {
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

resource "sysdig_monitor_alert_v2_group_outlier" "sample" {
	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
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

func alertV2GroupOutlierWithEnabled(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
	enabled = false

}
`, name)
}

func alertV2GroupOutlierWithUnreportedAlertNotificationsRetentionSec(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	metric = "sysdig_container_cpu_used_percent"
	group_aggregation = "avg"
	time_aggregation = "avg"
	group_by = ["kube_cluster_name", "cloud_provider_tag_Owner",]
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
  algorithm = "DBSCAN"
  dbscan_tolerance = 1
	observation_window_minutes = 15
	unreported_alert_notifications_retention_seconds = 60 * 60 * 24 * 30
}
`, name)
}
