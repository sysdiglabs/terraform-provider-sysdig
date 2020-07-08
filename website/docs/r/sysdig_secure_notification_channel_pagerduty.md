---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel_pagerduty"
sidebar_current: "docs-sysdig-secure-notification-channel-pagerduty"
description: |-
  Creates a Sysdig Secure Notification Channel of type Pagerduty.
---

# sysdig\_secure\_notification\_channel\_pagerduty

Creates a Sysdig Secure Notification Channel of type Pagerduty.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_notification_channel_pagerduty" "sample-pagerduty" {
	name                    = "Example Channel - Pagerduty"
	enabled                 = true
	account                 = "account"
	service_key             = "XXXXXXXXXX"
	service_name            = "sysdig"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `account` - (Required) Pagerduty account.

* `service_key` - (Required) Service Key for the Pagerduty account.

* `service_name` - (Required) Service name for the Pagerduty account.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.
