---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_promql"
description: |-
  Creates a Sysdig Monitor PromQL Alert.
---

# Resource: sysdig_monitor_alert_promql

Creates a Sysdig Monitor PromQL Alert. Monitor prometheus metrics and alert if they violate user-defined PromQL-based metric expression.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_monitor_alert_promql" "sample" {
	name = "Elasticsearch JVM heap usage"
	description = "A Kubernetes pod failed to restart"
	severity = 6

	promql = "(elasticsearch_jvm_memory_used_bytes{area=\"heap\"} / elasticsearch_jvm_memory_max_bytes{area=\"heap\"}) * 100 > 80"
	trigger_after_minutes = 10
}
```

## Argument Reference

### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `severity` - (Optional) Severity of the Monitor alert. It must be a value between 0 and 7,
               with 0 being the most critical and 7 the less critical. Defaults to 4.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Defaults to true.
* `notification_channels` - (Optional) List of notification channel IDs where an alert must be sent to once fired.
* `renotification_minutes` - (Optional) Number of minutes for the alert to re-notify until the status is solved.
* `custom_notification` - (Optional) Allows to define a custom notification title, prepend and append text.

### `custom_notification`

By defining this field, the user can modify the title and the body of the message sent when the alert
is fired.

* `title` - (Required) Sets the title of the alert. It is commonly defined as `{{__alert_name__}} is {{__alert_status__}}`.
* `prepend` - (Optional) Text to add before the alert template.
* `append` - (Optional) Text to add after the alert template.

### PromQL alert arguments

* `promql` - (Required) PromQL-based metric expression to alert on. Example: `histogram_quantile(0.99, rate(etcd_http_successful_duration_seconds_bucket[5m]) > 0.15` or `predict_linear(sysdig_fs_free_bytes{fstype!~"tmpfs"}[1h], 24*3600) < 10000000000`.

## Attributes Reference

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `id` - ID of the alert created.
* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.


## Import

PromQL Monitor alerts can be imported using the alert ID, e.g.

```
$ terraform import sysdig_monitor_alert_promql.example 12345
```