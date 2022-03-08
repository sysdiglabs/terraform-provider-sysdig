---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel"
description: |-
  Retrieves a Sysdig Secure Notification Channel.
---

# sysdig_secure_notification_channel

Retrieves the information of an existing Sysdig Secure Notification Channel.

~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_notification_channel" "sample-email" {
  name                 = "Example Channel - Email"
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `enabled` - If false, the channel will not emit notifications.

* `type` - Will be one of the following:  "EMAIL", "SNS", "OPSGENIE", 
    "VICTOROPS", "WEBHOOK", "SLACK", "PAGER_DUTY".

* `notify_when_ok` - Send a new notification when the alert condition is 
    no longer triggered.

* `notify_when_resolved` - Send a new notification when the alert is manually 
    acknowledged by a user.

* `send_test_notification` - Send an initial test notification to check
    if the notification channel is working.

### Attributes for type EMAIL

* `recipients` - Comma-separated list of recipients that will receive 
    the message.
    
### Attributes for type Amazon SNS

* `topics` - List of ARNs from the SNS topics.

### Attributes for type VICTOROPS

* `api_key` - Key for the API.

* `routing_key` - Routing key for VictorOps. 

### Attributes for type OPSGENIE

* `api_key` - Key for the API.

### Attributes for type WEBHOOK

* `url` - URL to send the event.

### Attributes for type SLACK

* `url` - URL of the Slack.

* `channel` - Channel name from this Slack.

### Attributes for type PAGERDUTY

* `account` - Pagerduty account.

* `service_key` - Service Key for the Pagerduty account.

* `service_name` - Service name for the Pagerduty account.
