---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_ibm_event_notification"
description: |-
  Retrieves information about a Monitor notification channel of type IBM Event Notification
---

# Data Source: sysdig_monitor_notification_channel_ibm_event_notification

Retrieves information about a Monitor notification channel of type IBM Event Notification.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_monitor_notification_channel_ibm_event_notification" "nc_ibm_event_notification" {
	name = "some notification channel name"
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Notification Channel ID.
* `name` - The Notification Channel Name.
* `instance_id` - id of the Event Notifications Instance. Id value can be either an instance id or CRN. If the event notification instance is within the same account, use the acutal instance id. If it is in a different accout, then use the Event Notifications Instance's [CRN](https://cloud.ibm.com/docs/account?topic=account-crn).
* `enabled` - Whether the Notification Channel is active or not.
* `notify_when_ok` - Whether the Notification Channel sends a notification when the condition is no longer triggered.
* `notify_when_resolved` - Whether the Notification Channel sends a notification if it's manually acknowledged by a
  user.
* `version` - The version of the Notification Channel.
* `send_test_notification` - Whether the Notification Channel has enabled the test notification.
