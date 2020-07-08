---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel_sns"
sidebar_current: "docs-sysdig-secure-notification-channel-sns"
description: |-
  Creates a Sysdig Secure Notification Channel of type Amazon SNS.
---

# sysdig\_secure\_notification\_channel\_sns

Creates a Sysdig Secure Notification Channel of type Amazon SNS.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_notification_channel_sns" "sample-amazon-sns" {
	name                    = "Example Channel - Amazon SNS"
	enabled                 = true
	topics                  = ["arn:aws:sns:us-east-1:273489009834:my-alerts2", "arn:aws:sns:us-east-1:279948934544:my-alerts"]
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `topics` - (Required) List of ARNs from the SNS topics.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.
