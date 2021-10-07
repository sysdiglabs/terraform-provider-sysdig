---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_event"
description: |-
  Creates a Sysdig Monitor Event Alert.
---

# Resource: sysdig\_monitor\_alert\_event

Creates a Sysdig Monitor Event Alert. Monitor occurrences of specific events, and alert if the total 
number of occurrences violates a threshold. Useful for alerting on container, orchestration, and 
service events like restarts and deployments.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_monitor_alert_event" "sample" {
	name = "[Kubernetes] Failed to pull image"
	description = "A Kubernetes pod failed to pull an image from the registry"
	severity = 4

	event_name = "Failed to pull image"
	source = "kubernetes"
	event_rel = ">"
	event_count = 0

	multiple_alerts_by = ["kubernetes.pod.name"]
	
	trigger_after_minutes = 1
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

### Event alert arguments

* `event_name` - (Required) String that matches part of name, tag or the description of Sysdig Events.
* `source` - (Required) Source of the event. It can be `docker` or `kubernetes`. 
* `event_rel` - (Required) Relationship of the event count. It can be `>`, `>=`, `<`, `<=`, `=` or `!=`.
* `event_count` - (Required) Number of events to match with event_rel.
* `multiple_alerts_by` - (Optional) List of segments to trigger a separate alert on. Example: `["kubernetes.cluster.name", "kubernetes.namespace.name"]`.  

## Attributes Reference

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `id` - ID of the alert created.
* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.


## Import

Event Monitor alerts can be imported using the alert ID, e.g.

```
$ terraform import sysdig_monitor_alert_event.example 12345
```