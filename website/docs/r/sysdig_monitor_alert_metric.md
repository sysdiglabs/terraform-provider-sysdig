---
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_alert_metric"
sidebar_current: "docs-sysdig-monitor-alert-metric"
description: |-
  Creates a Sysdig Monitor Metric Alert.
---

# sysdig\_monitor\_alert\_metric

Creates a Sysdig Monitor Metric Alert. Monitor time-series metrics and alert if they violate user-defined thresholds.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_monitor_alert_metric" "sample" {
	name = "[Kubernetes] CrashLoopBackOff"
	description = "A Kubernetes pod failed to restart"
	severity = 6

	metric = "sum(timeAvg(kubernetes.pod.restart.count)) > 2"
	trigger_after_minutes = 1

	multiple_alerts_by = ["kubernetes.cluster.name",
                          "kubernetes.namespace.name",
                          "kubernetes.deployment.name",
                          "kubernetes.pod.name"]

	capture {
		filename = "CrashLoopBackOff"
		duration = 15
	}
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

* `metric` - (Required) Metric to monitor and alert on. Example: `sum(timeAvg(kubernetes.pod.restart.count)) > 2` or `avg(avg(cpu.used.percent)) > 50`.
* `multiple_alerts_by` - (Optional) List of segments to trigger a separate alert on. Example: `["kubernetes.cluster.name", "kubernetes.namespace.name"]`.  

## Attributes Reference

### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.