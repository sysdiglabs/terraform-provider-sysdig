---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_email"
description: |-
  Creates a Sysdig Monitor Notification Channel of type Email.
---

# Resource: sysdig_monitor_notification_channel_email

Creates a Sysdig Monitor Notification Channel of type Email.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_notification_channel_email" "sample_email" {
	name                    = "Example Channel - Email"
	recipients              = ["foo@localhost.com", "bar@localhost.com"]
	enabled                 = true
	send_test_notification  = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `recipients` - (Required) List of recipients that will receive
    the message.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional, Deprecated) Send a new notification when the alert condition is no longer triggered. Default is `false`. This option is deprecated; use `notify_on_resolve` within the `notification_channels` options in the `sysdig_monitor_alert_v2_*` resources instead, which takes precedence over this setting.

* `notify_when_resolved` - (Optional, Deprecated) Send a new notification when the alert is manually acknowledged by a user. Default is `false`. This option is deprecated; use `notify_on_acknowledge` within the `notification_channels` options in the `sysdig_monitor_alert_v2_*` resources instead, which takes precedence over this setting.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.

* `share_with_current_team` - (Optional) If set to `true` it will share notification channel only with current team (in which user is logged in).
  Otherwise, it will share it with all teams, which is the default behaviour. Although this is an optional setting, beware that if you have lower permissions than admin you may see a `error: 403 Forbidden` if this is not set to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the Notification Channel.

* `version` - (Computed) The current version of the Notification Channel.

## Import

Email notification channels for Monitor can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_notification_channel_email.example 12345
```
