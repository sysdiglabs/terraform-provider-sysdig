---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_v2_prometheus"
description: |-
  Creates a Sysdig Monitor PromQL Alert with AlertV2 API.
---

# Resource: sysdig_monitor_alert_v2_prometheus

Monitor your infrastructure with PromQL queries, maintaining full compatibility with OSS Prometheus.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_alert_v2_prometheus" "sample" {
  name = "Elasticsearch JVM heap usage"
  description = "Elasticsearch JVM heap used over attention threshold"
  severity = "high"
  query = "(elasticsearch_jvm_memory_used_bytes{area=\"heap\"} / elasticsearch_jvm_memory_max_bytes{area=\"heap\"}) * 100 > 80"
  duration_seconds = 600
  notification_channels {
    id = 1234
    renotify_every_minutes = 5
  }
  labels = {
    application = "app1"
    maturity = "high"
  }
}
```

## Argument Reference

### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the alert rule. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `duration_seconds` - (Optional) Specifies the amount of time, in seconds, that an alert condition must remain continuously true before the alert rule is triggered.
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
* `notify_on_resolve` - (Optional) Whether to send a notification when the alert is resolved. Default: `true`.
* `notify_on_acknowledge` - (Optional) Whether to send a notification when the alert is acknowledged. If not defined, this option is inherited from the `notify_when_resolved` option from the specific notification channel selected.

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

### Prometheus alert arguments

* `query` - (Required) PromQL-based metric expression to alert on. Example: `histogram_quantile(0.99, rate(etcd_http_successful_duration_seconds_bucket[5m]) > 0.15` or `predict_linear(sysdig_fs_free_bytes{fstype!~"tmpfs"}[1h], 24*3600) < 10000000000`.
* `keep_firing_for_minutes` - (Optional) Alert resolution delay before actually resolving an alert.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the alerts in Sysdig Monitor:

* `id` - ID of the alert created.
* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.

## Import

Prometheus alerts can be imported using the alert ID, e.g.

```
$ terraform import sysdig_monitor_alert_v2_prometheus.example 12345
```
