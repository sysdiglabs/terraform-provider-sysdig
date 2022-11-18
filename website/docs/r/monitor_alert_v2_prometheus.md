---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_v2_prometheus"
description: |-
  Creates a Sysdig Monitor PromQL Alert with AlertV2 API.
---

# Resource: sysdig_monitor_alert_v2_prometheus

Creates a Sysdig Monitor Prometheus Alert. The notification is triggered on the user-defined PromQL expression.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_alert_v2_prometheus" "sample" {
	name = "Elasticsearch JVM heap usage"
	description = "A Kubernetes pod failed to restart"
	severity = high
	query = "(elasticsearch_jvm_memory_used_bytes{area=\"heap\"} / elasticsearch_jvm_memory_max_bytes{area=\"heap\"}) * 100 > 80"
	trigger_after_minutes = 10
	notification_channels {
		id = <your-notification-channel-id>
    type = "EMAIL""
    renotify_every_minutes = 5
	}
}
```

## Argument Reference

### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `group` - (Optional) Lowercase string to group alerts in the UI
* `severity` - (Optional) Severity of the Monitor alert. It must be `high`, `medium`, `low` or `info`. Default: `low`.
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Defaults to true.
* `notification_channels` - (Optional) List of notification channel configuration
* `custom_notification` - (Optional) Allows to define a custom notification title, prepend and append text.
* `capture` - (Optional) Allows to define a configuration to trigger a Sysdig Capture.

### `notification_channels` - 

By defining this field, the user can choose to which notification channels send the events when the alert fires. 

It is a list of objects with the following fields:
* `id` - (Required) The ID of the notification channel
* `type` - (Required) The type of the notification channel
* `renotify_every_minutes`: (Optional) the amount of minutes to wait before re sending the notification to this channel

### `custom_notification`

By defining this field, the user can modify the title and the body of the message sent when the alert
is fired.

* `subject` - (Optional) Sets the title of the alert.
* `prepend` - (Optional) Text to add before the alert template.
* `append` - (Optional) Text to add after the alert template.

### Prometheus alert arguments

* `query` - (Required) PromQL-based metric expression to alert on. Example: `histogram_quantile(0.99, rate(etcd_http_successful_duration_seconds_bucket[5m]) > 0.15` or `predict_linear(sysdig_fs_free_bytes{fstype!~"tmpfs"}[1h], 24*3600) < 10000000000`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `id` - ID of the alert created.
* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.


## Import

Prometheus Monitor alerts can be imported using the alert ID, e.g.

```
$ terraform import sysdig_monitor_alert_v2_prometheus.example 12345
```