---
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_email"
sidebar_current: "docs-sysdig-monitor-notification-channel-email"
description: |-
  Creates a Sysdig Monitor Notification Channel of type Email.
---

# sysdig\_monitor\_notification_channel\_email

Creates a Sysdig Monitor Notification Channel of type Email.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_monitor_notification_channel_email" "sample_email" {
	name                    = "Example Channel - Email"
	recipients              = ["foo@localhost.com", "bar@localhost.com"]
	enabled                 = true
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `recipients` - (Required) List of recipients that will receive 
    the message.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.
