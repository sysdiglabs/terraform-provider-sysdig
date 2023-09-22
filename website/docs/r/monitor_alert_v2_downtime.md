---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_v2_downtime"
description: |-
  Creates a Sysdig Monitor Downtime Alert with AlertV2 API.
---

# Resource: sysdig_monitor_alert_v2_downtime

Creates a Sysdig Monitor Downtime Alert. Monitor any type of entity - host, container, process - and alert when the entity goes down.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_alert_v2_downtime" "sample" {

  name = "process down"
  severity = "high"
  metric = "sysdig_program_up"
  threshold = 75
  group_by = ["host_hostname", "program_name"]

  scope {
    label = "host_hostname"
    operator = "in"
    values = ["my-cluster-1", "my-server-2"]
  }

  notification_channels {
    id = 1234
  }

  trigger_after_minutes = 10

}

```

## Argument Reference

### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `group` - (Optional) Lowercase string to group alerts in the UI.
* `severity` - (Optional) Severity of the Monitor alert. It must be `high`, `medium`, `low` or `info`. Default: `low`.
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Default: `true`.
* `notification_channels` - (Optional) List of notification channel configurations.
* `custom_notification` - (Optional) Allows to define a custom notification title, prepend and append text.
* `link` - (Optional) List of links to add to notifications.

### `notification_channels`

By defining this field, the user can choose to which notification channels send the events when the alert fires.

It is a list of objects with the following fields:
* `id` - (Required) The ID of the notification channel.
* `renotify_every_minutes` - (Optional) the amount of minutes to wait before re sending the notification to this channel. `0` means no renotification enabled. Default: `0`.
* `notify_on_resolve` - (Optional) Wether to send a notification when the alert is resolved. Default: `true`.

### `custom_notification`

By defining this field, the user can modify the title and the body of the message sent when the alert is fired.

* `subject` - (Optional) Sets the title of the alert.
* `prepend` - (Optional) Text to add before the alert template.
* `append` - (Optional) Text to add after the alert template.

### `link`

By defining this field, the user can add link to notifications.

* `type` - (Required) Type of link. Must be `runbook`, for generic links, or `dashboard`, for internal links to existing dashboards.
* `href` - (Optional) When using `runbook` type, url of the external resource.
* `id` - (Optional) When using `dashboard` type, dashboard id.

### `capture`

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file. Must have `.scap` suffix.
* `duration_seconds` - (Optional) Time frame of the capture. Default: `15`.
* `storage` - (Optional) Custom bucket where to save the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.
* `enabled` - (Optional) Wether to enable captures. Default: `true`.

### Downtime alert arguments

* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure. Can be repeated.
* `group_by` - (Optional) List of segments to trigger a separate alert on. Example: `["kube_cluster_name", "kube_pod_name"]`.
* `metric` - (Required) Metric the alert will act upon. Can be: `sysdig_container_up`, `sysdig_program_up`, `sysdig_host_up`.
* `threshold` - (Required) Below of this percentage of downtime the alert will be triggered. Defaults to 100.
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

Downtime alerts can be imported using the alert ID, e.g.

```
$ terraform import sysdig_monitor_alert_v2_downtime.example 12345
```
