---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_sns"
description: |-
  Creates a Sysdig Monitor Notification Channel of type Amazon SNS.
---

# Resource: sysdig_monitor_notification_channel_sns

Creates a Sysdig Monitor Notification Channel of type Amazon SNS.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example Usage

```hcl
resource "sysdig_monitor_notification_channel_sns" "sample-amazon-sns" {
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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the Notification Channel.

* `version` - (Computed) The current version of the Notification Channel.

## Import

Amazon SNS notification channels for Monitor can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_notification_channel_sns.example 12345
```