---
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_anomaly"
sidebar_current: "docs-sysdig-monitor-alert-anomaly"
description: |-
  Creates a Sysdig Monitor Anomaly Alert.
---

# sysdig\_monitor\_alert\_anomaly

Creates a Sysdig Monitor Anomaly Alert. Monitor hosts based on their historical behaviors and alert when they deviate.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_monitor_alert_anomaly" "sample" {
	name = "[Kubernetes] Anomaly Detection Alert"
	description = "Detects an anomaly in the cluster"
	severity = 6

	monitor = ["cpu.used.percent", "memory.bytes.used"]

	multiple_alerts_by = ["kubernetes.cluster.name", 
                          "kubernetes.namespace.name", 
                          "kubernetes.deployment.name", 
                          "kubernetes.pod.name"]
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
 

#### Capture

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file.
* `duration` - (Required) Time frame in seconds of the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.

### Metric alert arguments

* `monitor` - (Required) Array of metrics to monitor and alert on. Example: `["cpu.used.percent", "cpu.cores.used", "memory.bytes.used", "fs.used.percent", "thread.count", "net.request.count.in"]`.
* `multiple_alerts_by` - (Optional) List of segments to trigger a separate alert on. Example: `["kubernetes.cluster.name", "kubernetes.namespace.name"]`.  

## Attributes Reference

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.