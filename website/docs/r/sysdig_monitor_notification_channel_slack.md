---
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_slack"
sidebar_current: "docs-sysdig-monitor-notification-channel-slack"
description: |-
  Creates a Sysdig Monitor Notification Channel of type Slack.
---

# sysdig\_monitor\_notification\_channel\_slack

Creates a Sysdig Monitor Notification Channel of type Slack.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example usage

```hcl
resource "sysdig_monitor_notification_channel_slack" "sample-slack" {
	name                    = "Example Channel - Slack"
	enabled                 = true
	url                     = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel                 = "#sysdig"
	notify_when_ok          = false
	notify_when_resolved    = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `url` - (Required) URL of the Slack.

* `channel` - (Required) Channel name from this Slack.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.
