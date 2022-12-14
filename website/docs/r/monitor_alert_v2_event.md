---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_v2_event"
description: |-
  Creates a Sysdig Monitor Event Alert with AlertV2 API.
---

# Resource: sysdig_monitor_alert_v2_event

Creates a Sysdig Monitor Event Alert. Monitor occurrences of specific events, and alert if the total 
number of occurrences violates a threshold. Useful for alerting on container, orchestration, and 
service events like restarts and deployments.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_alert_v2_event" "sample" {

  name = "[Kubernetes] Failed to pull image"
  description = "A Kubernetes pod failed to pull an image from the registry"
  severity = "high"
  filter = "Failed to pull image"
  sources = ["kubernetes"]
  operator = ">"
  threshold = 0
  group_by = ["kube_pod_name"]

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

  trigger_after_minutes = 1

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
* `notification_channels` - (Optional) List of notification channel configuration
* `custom_notification` - (Optional) Allows to define a custom notification title, prepend and append text.
* `link` - (Optional) List of links to add to notifications.

### `notification_channels`

By defining this field, the user can choose to which notification channels send the events when the alert fires. 

It is a list of objects with the following fields:
* `id` - (Required) The ID of the notification channel.
* `renotify_every_minutes`: (Optional) the amount of minutes to wait before re sending the notification to this channel. `0` means no renotifiacation enabled. Default: `0`.
* `main_threshold`: (Optional) Whether this notification channel is used for the main threshold of the alert. Default: `true`.
* `warning_threshold`: (Optional) Whether this notification channel is used for the warning threshold of the alert. Default: `false`.

### `custom_notification`

By defining this field, the user can modify the title and the body of the message sent when the alert
is fired.

* `subject` - (Optional) Sets the title of the alert.
* `prepend` - (Optional) Text to add before the alert template.
* `append` - (Optional) Text to add after the alert template.

### `link`

By defining this field, the user can add link to notificatons.

* `type` - (Required) Type of link. Must be `runbook`, for generic links, or `dashboard`, for internal links to existing dashboards.
* `href` - (Optional) When using `runbook` type, url of the external resource.
* `id` - (Optional) When using `dashboard` type, dasboard id.

### `capture`

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file. Must have `.scap` suffix.
* `duration_seconds` - (Optional) Time frame of the capture. Default: `15`.
* `storage` - (Optional) Custom bucket where to save the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.
* `enabled` - (Optional) Wether to enable captures. Default: `true`.

### Event alert arguments

* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure. Can be repeated.
* `group_by` - (Optional) List of segments to trigger a separate alert on. Example: `["kube_cluster_name", "kube_pod_name"]`.
* `operator` - (Required) Condition operator of the event count. It can be `>`, `>=`, `<`, `<=`, `=` or `!=`.
* `threshold` - (Required) Number of events to match with `op`.
* `warning_threshold` - (Optional) Number of events to match with `op` to trigger a warning alert. Must be a number lower than `threshold`.
* `filter` - (Required) String that matches part of name, tag or the description of Sysdig Events.
* `sources` - (Required) List of sources of the event. It can be `kubernetes`, `containerd`, `docker` or arbitrary custom sources.

### `scope`

* `label` - (Required) Label in prometheus notation to select a part of the infrastructure.
* `operator` - (Required) Operator to match the label. It can be `equals`, `notEquals`, `in`, `notIn`, `contains`, `notContains`, `startsWith`.
* `values` - (Required) List of values to match the scope.

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
$ terraform import sysdig_monitor_alert_v2_event.example 12345
```
