---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel"
sidebar_current: "docs-sysdig-secure-notification-channel"
description: |-
  Creates a Sysdig Secure Notification Channel.
---

# sysdig\_secure\_notification_channel

Creates a Sysdig Secure Notification Channel.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_notification_channel" "sample-email" {
  name                 = "Example Channel - Email"
  enabled              = true
  type                 = "EMAIL"
  recipients           = "root@localhost.com"
  notify_when_ok       = false
  notify_when_resolved = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `enabled` - (Required) If false, the channel will not emit notifications.

* `type` - (Required) Must be one of the following:  "EMAIL", "SNS", "OPSGENIE", 
    "VICTOROPS", "WEBHOOK", "SLACK", "PAGER_DUTY".

* `notify_when_ok` - (Required) Send a new notification when the alert condition is 
    no longer triggered.

* `notify_when_resolved` - (Required) Send a new notification when the alert is manually 
    acknowledged by a user.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working.

### Arguments for type EMAIL

* `recipients` - (Required) Comma-separated list of recipients that will receive 
    the message.
    
### Arguments for type Amazon SNS

* `topics` - (Required) List of ARNs from the SNS topics.

### Arguments for type VICTOROPS

* `api_key` - (Required) Key for the API.

* `routing_key` - (Required) Routing key for VictorOps. 

### Arguments for type OPSGENIE

* `api_key` - (Required) Key for the API.

### Arguments for type WEBHOOK

* `url` - (Required) URL to send the event.

### Arguments for type SLACK

* `url` - (Required) URL of the Slack.

* `channel` - (Required) Channel name from this Slack.

### Arguments for type PAGERDUTY

* `account` - (Required) Pagerduty account.

* `service_key` - (Required) Service Key for the Pagerduty account.

* `service_name` - (Required) Service name for the Pagerduty account.
