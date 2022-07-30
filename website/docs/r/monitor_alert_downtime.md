---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_downtime"
description: |-
  Creates a Sysdig Monitor Downtime Alert.
---

# Resource: sysdig_monitor_alert_downtime

Creates a Sysdig Monitor Downtime Alert. Monitor any type of entity - host, container, process, service, etc - and alert when the entity goes down.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_alert_downtime" "sample" {
	name = "[Kubernetes] Downtime Alert"
	description = "Detects a downtime in the Kubernetes cluster"
	severity = 2

	entities_to_monitor = ["kubernetes.namespace.name"]

	trigger_after_minutes = 10
	trigger_after_pct = 100
}
```

## Argument Reference

### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `group` - (Optional) The group of Monitor alert.
* `severity` - (Optional) Severity of the Monitor alert. It must be a value between 0 and 7,
               with 0 being the most critical and 7 the less critical. Defaults to 4.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure.
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Defaults to true.
* `notification_channels` - (Optional) List of notification channel IDs where an alert must be sent to once fired.
* `renotification_minutes` - (Optional) Number of minutes for the alert to re-notify until the status is solved.
* `capture` - (Optional) Enables the creation of a capture file of the syscalls during the event.
* `custom_notification` - (Optional) Allows to define a custom notification title, prepend and append text.

### `capture`

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file.
* `duration` - (Required) Time frame in seconds of the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.

### `custom_notification`

By defining this field, the user can modify the title and the body of the message sent when the alert
is fired.

* `title` - (Required) Sets the title of the alert. It is commonly defined as `{{__alert_name__}} is {{__alert_status__}}`.
* `prepend` - (Optional) Text to add before the alert template.
* `append` - (Optional) Text to add after the alert template.

### Downtime alert arguments

* `entities_to_monitor` - (Required) List of metrics to monitor downtime and alert on. Example: `["kubernetes.namespace.name"]` to detect namespace removal or `["host.hostName"]` to detect host downtime.
* `trigger_after_pct` - (Optional) Below of this percentage of downtime the alert will be triggered. Defaults to 100.  

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `id` - ID of the alert created.
* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.


## Import

Downtime Monitor alerts can be imported using the alert ID, e.g.

```
$ terraform import sysdig_monitor_alert_downtime.example 12345
```
