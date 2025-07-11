---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_v2_form_based_prometheus"
description: |-
  Creates a Sysdig Monitor Threshold Prometheus Alert with AlertV2 API.
---

# Resource: sysdig_monitor_alert_v2_form_based_prometheus

-> **Note:** Form Based Prometheus Alerts are now part of Threshold Alerts. The Terraform resource remains `sysdig_monitor_alert_v2_form_based_prometheus` for backwards compatibility.

Threshold Alerts configured with PromQL allow you to monitor your infrastructure by comparing any PromQL expression against user-defined thresholds.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_alert_v2_form_based_prometheus" "sample" {
  name = "Elasticsearch JVM heap usage"
  description = "Elasticsearch JVM heap used over attention threshold"
  severity = "high"
  query = "(elasticsearch_jvm_memory_used_bytes{area=\"heap\"} / elasticsearch_jvm_memory_max_bytes{area=\"heap\"}) * 100"
  operator = ">"
  threshold = 80
  notification_channels {
    id = 1234
    renotify_every_minutes = 5
  }
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
* `renotify_every_minutes` - (Optional) the amount of minutes to wait before re sending the notification to this channel. `0` means no renotification enabled.
* `notify_on_resolve` - (Optional) Wether to send a notification when the alert is resolved. Default: `true`.
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

* `type` - (Required) Type of link. Must be `runbook`, for generic links, or `dashboard`, for internal links to existing dashboards.
* `href` - (Optional) When using `runbook` type, url of the external resource.
* `id` - (Optional) When using `dashboard` type, dashboard id.

### Threshold Prometheus alert arguments

* `query` - (Required) PromQL-based metric expression to alert on. Example: `sysdig_host_memory_available_bytes / sysdig_host_memory_total_bytes * 100` or `avg_over_time(sysdig_container_cpu_used_percent{}[59s])`.
* `operator` - (Required) Operator for the condition to alert on. It can be `>`, `>=`, `<`, `<=`, `==` or `!=`.
* `threshold` - (Required) Threshold used together with `op` to trigger the alert if crossed.
* `warning_threshold` - (Optional) Warning threshold used together with `op` to trigger the alert if crossed. Must be a number that triggers the alert before reaching the main `threshold`.
* `duration_seconds` - (Optional) Specifies the amount of time, in seconds, that an alert condition must remain continuously true before the alert rule is triggered.
* `no_data_behaviour` - (Optional) behaviour in case of missing data. Can be `DO_NOTHING`, i.e. ignore, or `TRIGGER`, i.e. notify on main threshold. Default: `DO_NOTHING`.
* `unreported_alert_notifications_retention_seconds` - (Optional) Period after which any alerts triggered for entities (such as containers or hosts) that are no longer reporting data will be automatically marked as 'deactivated'. By default there is no deactivation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the alerts in Sysdig Monitor:

* `id` - ID of the alert created.
* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.

## Import

Threshold Prometheus alerts can be imported using the alert ID, e.g.

```
$ terraform import sysdig_monitor_alert_v2_form_based_prometheus.example 12345
```
