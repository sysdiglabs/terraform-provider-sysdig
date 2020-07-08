---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel_victorops"
sidebar_current: "docs-sysdig-secure-notification-channel-victorops"
description: |-
  Creates a Sysdig Secure Notification Channel of type VictorOps.
---

# sysdig\_secure\_notification\_channel\_victorops

Creates a Sysdig Secure Notification Channel of type VictorOps.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_notification_channel_victorops" "sample-victorops" {
	name                    = "Example Channel - VictorOps"
	enabled                 = true
	api_key                 = "1234342-4234243-4234-2"
	routing_key             = "My team"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `api_key` - (Required) Key for the API.

* `routing_key` - (Required) Routing key for VictorOps. 

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.
