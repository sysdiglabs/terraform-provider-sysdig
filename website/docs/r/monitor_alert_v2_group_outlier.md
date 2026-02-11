---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_v2_group_outlier"
description: |-
  Creates a Sysdig Monitor Group Outlier Alert with AlertV2 API.
---

# Resource: sysdig_monitor_alert_v2_group_outlier

Creates a Sysdig Monitor Group Outlier Alert. Monitor specific segments in a metric to identify entities that deviate from the group.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_alert_v2_group_outlier" "sample" {

  name = "cpu usage outliers"
  severity = "high"
  metric = "sysdig_container_cpu_used_percent"

  algorithm = "MAD"
  mad_threshold = 10.1
  mad_tolerance = 5.5

  group_aggregation = "avg"
  group_by = ["kube_pod_name", "container_name"]
  time_aggregation = "avg"

  scope {
    label = "kube_cluster_name"
    operator = "in"
    values = ["my_cluster_1", "my_cluster_2"]
  }

  scope {
    label = "kube_deployment_name"
    operator = "equals"
    values = ["my_deployment"]
  }

  notification_channels {
    id = 1234
    renotify_every_minutes = 60
  }

  observation_window_minutes = 15

}

```

## Argument Reference

### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the alert rule. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `group` - (Optional) Used to group alert rules in the UI. This value must be a lowercase string.
* `severity` - (Optional) Severity of the Monitor alert. It must be `high`, `medium`, `low` or `info`. Default: `low`.
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Default: `true`.
* `notification_channels` - (Optional) List of notification channel configurations.
* `custom_notification` - (Optional) Allows to define a custom notification title, prepend and append text.
* `link` - (Optional) List of links to add to notifications.
* `labels` - (Optional) map of labels to be attached to this alert.

### `notification_channels`

By defining this field, the user can choose to which notification channels send the events when the alert fires.

It is a list of objects with the following fields:
* `id` - (Required) The ID of the notification channel.
* `renotify_every_minutes` - (Optional) the amount of minutes to wait before re sending the notification to this channel. `0` means no renotification enabled. Default: `0`.
* `notify_on_resolve` - (Optional) Whether to send a notification when the alert is resolved. Default: `true`.
* `notify_on_acknowledge` - (Optional) Whether to send a notification when the alert is acknowledged. If not defined, this option is inherited from the `notify_when_resolved` option from the specific notification channel selected.
* `main_threshold` - (Optional) Whether this notification channel is used for the main threshold of the alert. Default: `true`.
* `warning_threshold` - (Optional) Whether this notification channel is used for the warning threshold of the alert. Default: `false`.

### `custom_notification`

By defining this field, the user can modify the title and the body of the message sent when the alert is fired.

* `subject` - (Optional) Sets the title of the alert.
* `prepend` - (Optional) Text to add before the alert template.
* `append` - (Optional) Text to add after the alert template.
* `additional_field` - (Optional) Set of additional fields to add to the notification.

#### `additional_field`
* `name` - (Required) field name.
* `value` - (Required) field value.

### `link`

By defining this field, the user can add link to notifications.

* `type` - (Required) Type of link. Must be `runbook` for generic links, `dashboard` for internal links to existing dashboards, or `dashboardTemplate` for links to dashboard templates.
* `href` - (Optional) When using `runbook` type, url of the external resource.
* `id` - (Optional) When using `dashboard` type, dashboard id. When using `dashboardTemplate` type, the dashboard template id (e.g. `view.promcat.mysql`).

### `capture`

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file. Must have `.scap` suffix.
* `duration_seconds` - (Optional) Time frame of the capture. Default: `15`.
* `storage` - (Optional) Custom bucket where to save the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.
* `enabled` - (Optional) Whether to enable captures. Default: `true`.

### Group Outlier alert arguments

* `observation_window_minutes` - (Required) Specific time frame in minutes for evaluating potential outliers. The minimum value is ten minutes.
* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure. Can be repeated.
* `group_by` - (Required) List of segments to trigger a separate alert on. Example: `["kube_cluster_name", "kube_pod_name"]`.
* `metric` - (Required) Metric the alert will act upon.
* `time_aggregation` - (Required) time aggregation function for data. It can be `avg`, `timeAvg`, `sum`, `min`, `max`.
* `group_aggregation` - (Required) group aggregation function for data. It can be `avg`, `sum`, `min`, `max`.
* `algorithm` - (Required) Algorithm to use to detect outliers. Can be `MAD` (Median Absolute Deviation) or `DBSCAN` (Density-Based Spatial Clustering of Applications with Noise).
* `dbscan_tolerance` - (Optional - Required if `algorithm = DBSCAN`) Proximity range within which an entity should find neighboring time series to be part of a group. Allowed values are between 0.5 and 10.
* `mad_tolerance` - (Optional - Required if `algorithm = MAD`) Tolerance to decide the acceptable values from the median absolute deviation. Allowed values are between 0.5 and 10.
* `mad_threshold` - (Optional - required if `algorithm = MAD`) Percentage of the observation window in which an entity’s reported value must fall outside the configured tolerance to be labeled as an outlier. Allowed values are between 1 and 100.
* `no_data_behaviour` - (Optional) behaviour in case of missing data. Can be `DO_NOTHING`, i.e. ignore, or `TRIGGER`, i.e. notify on main threshold. Default: `DO_NOTHING`.
* `unreported_alert_notifications_retention_seconds` - (Optional) Period after which any alerts triggered for entities (such as containers or hosts) that are no longer reporting data will be automatically marked as 'deactivated'. By default there is no deactivation.

### `scope`

* `label` - (Required) Label in prometheus notation to select a part of the infrastructure.
* `operator` - (Required) Operator to match the label. It can be `equals`, `notEquals`, `in`, `notIn`, `contains`, `notContains`, `startsWith`.
* `values` - (Required) List of values to match the scope.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the alerts in Sysdig Monitor:

* `id` - ID of the alert created.
* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.


## Import

Group Outlier alerts can be imported using the alert ID, e.g.

```
$ terraform import sysdig_monitor_alert_v2_group_outlier.example 12345
```
