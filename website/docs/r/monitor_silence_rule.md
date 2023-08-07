---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_silence_rule"
description: |-
  Creates a Sysdig Monitor Silence Rule.
---

# Resource: sysdig_monitor_silence_rule

Creates a Sysdig Monitor Silence Rule.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "time_static" "start_ts" {
  rfc3339 = "2023-07-08T07:00:00Z"
}

resource "sysdig_monitor_silence_rule" "sample" {
  name = "Example Silence Rule"
  enabled = true
  start_ts = time_static.start_ts.unix * 1000
  duration_seconds = 60 * 60 * 24
  scope = "cloudProvider.region != \"us-east-1\" and not host.hostName contains \"testhost\" and kubernetes.job.name starts with \"prod\" and kubernetes.daemonSet.name in (\"ds1\", \"ds2\")"
  alert_ids = [1234, 1235]
  notification_channel_ids = [111, 222]
}
```

## Argument Reference

Ended Silence Rules cannot be updated.

* `name` - (Required) The name of the Silence Rule.

* `enabled` - (Optional) Whether to enable the Silence Rule. Default: `true`.

* `start_ts` - (Required) Unix timestamp, in milliseconds, when the Silence Rule starts.

* `duration_seconds` - (Required) Duration of the Silence Rule, in seconds.

* `scope` - (Optional) Part of the infrastructure the Silence Rule will be applied to. At least one of `scope` or `alert_ids` must be defined.

* `alert_ids` - (Optional) List of alerts the Silence Rule will be applied to. At least one of `scope` or `alert_ids` must be defined.

* `notification_channel_ids` - (Optional) List of notification channels that will be used to notify when the Silence Rule starts and end.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the Silence Rule.

* `version` - (Computed) The current version of the Silence Rule.

## Import

Silence Rules for Monitor can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_silence_rule.example 12345
```
