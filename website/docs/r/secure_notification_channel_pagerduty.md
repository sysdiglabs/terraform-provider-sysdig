---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel_pagerduty"
description: |-
  Creates a Sysdig Secure Notification Channel of type Pagerduty.
---

# Resource: sysdig_secure_notification_channel_pagerduty

Creates a Sysdig Secure Notification Channel of type Pagerduty.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_notification_channel_pagerduty" "sample-pagerduty" {
	name                    = "Example Channel - Pagerduty"
	enabled                 = true
	account                 = "account"
	service_key             = "XXXXXXXXXX"
	service_name            = "sysdig"
	send_test_notification  = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `account` - (Required) Pagerduty account.

* `service_key` - (Required) Service Key for the Pagerduty account.

* `service_name` - (Required) Service name for the Pagerduty account.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional, Deprecated) Send a new notification when the alert condition is no longer triggered. Default is `true`. This option is deprecated; use `notify_on_resolve` within the `notification_channels` options in the `sysdig_monitor_alert_v2_*` resources instead, which takes precedence over this setting. This option only applies to Monitor alerts when the channel is shared across all teams. It has no effect on Secure features.

* `notify_when_resolved` - (Optional, Deprecated) Send a new notification when the alert is manually acknowledged by a user. Default is `true`. This option is deprecated; use `notify_on_acknowledge` within the `notification_channels` options in the `sysdig_monitor_alert_v2_*` resources instead, which takes precedence over this setting. This option only applies to Monitor alerts when the channel is shared across all teams. It has no effect on Secure features.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the Notification Channel.

* `version` - (Computed) The current version of the Notification Channel.

* `share_with_current_team` - (Optional) If set to `true` it will share notification channel only with current team (in which user is logged in).
  Otherwise, it will share it with all teams, which is the default behaviour. Although this is an optional setting, beware that if you have lower permissions than admin you may see a `error: 403 Forbidden` if this is not set to `true`.

## Import

Pagerduty notification channels for Secure can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_notification_channel_pagerduty.example 12345
```
