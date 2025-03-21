---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_ibm_event_notification"
description: |-
  Creates a Sysdig Monitor Notification Channel of type IBM Event Notification.
---

# Resource: sysdig_monitor_notification_channel_ibm_event_notification

Creates a Sysdig Monitor Notification Channel of type IBM Event Notification (only available in IBM Cloud Monitoring).

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
# IBM event notification registering in the same account
resource "sysdig_monitor_notification_channel_ibm_event_notification" "sample" {
	name                    = "Example Channel - IBM Event Notification"
	enabled                 = true
	instance_id             = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	notify_when_ok          = false
	notify_when_resolved    = false
	share_with_current_team = true
}
```

```terraform
# IBM event notification registering in different account
resource "sysdig_monitor_notification_channel_ibm_event_notification" "sample" {
	name                    = "Example Channel - IBM Event Notification"
	enabled                 = true
	instance_id             = "crn:v1:bluemix:public:event-notifications:global:a/59bcbfa6ea2f006b4ed7094c1a08dcdd:1a0ec336-f391-4091-a6fb-5e084a4c56f4::"
	notify_when_ok          = false
	notify_when_resolved    = false
	share_with_current_team = true
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `instance_id` - (Required) id of the Event Notifications Instance. Id value can be either an instance id or CRN. If the event notification instance is within the same account, use the actual instance id. If it is in a different account, then use the Event Notifications Instance's [CRN](https://cloud.ibm.com/docs/account?topic=account-crn).

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.

* `share_with_current_team` - (Optional) If set to `true` it will share notification channel only with current team (in which user is logged in).
  Otherwise, it will share it with all teams, which is the default behaviour. Although this is an optional setting, beware that if you have lower permissions than admin you may see a `error: 403 Forbidden` if this is not set to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the Notification Channel.

* `version` - (Computed) The current version of the Notification Channel.

## Import

IBM Event Notification notification channels for Monitor can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_notification_channel_ibm_event_notification.example 12345
```
