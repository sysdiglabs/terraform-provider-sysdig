---
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_webhook"
sidebar_current: "docs-sysdig-monitor-notification-channel-webhook"
description: |-
  Creates a Sysdig Monitor Notification Channel of type Webhook.
---

# sysdig\_monitor\_notification\_channel\_webhook

Creates a Sysdig Monitor Notification Channel of type Webhook.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_monitor_notification_channel_webhook" "sample-webhook" {
	name                    = "Example Channel - Webhook"
	enabled                 = true
	url                     = "localhost:8080"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `url` - (Required) URL to send the event.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.
