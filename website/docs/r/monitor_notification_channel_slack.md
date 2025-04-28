---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_slack"
description: |-
  Creates a Sysdig Monitor Notification Channel of type Slack.
---

# Resource: sysdig_monitor_notification_channel_slack

Creates a Sysdig Monitor Notification Channel of type Slack.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_notification_channel_slack" "sample-slack" {
	name                    = "Example Channel - Slack"
	enabled                 = true
	url                     = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel                 = "#sysdig"
	is_private_channel      = false
	notify_when_ok          = false
	notify_when_resolved    = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `url` - (Required) URL of the Slack.

* `show_section_runbook_links` - (Optional) Whether to include the runbook links section in the Slack messages. Default: true.

* `show_section_event_details` - (Optional) Whether to include the event details section in the Slack messages. Default: true.

* `show_section_user_defined_content` - (Optional) Whether to include the user defined section in the Slack messages. Default: true.

* `show_section_notification_chart` - (Optional) Whether to include the notification chart section in the Slack messages. Default: true.

* `show_section_dashboard_links` - (Optional) Whether to include the dashboard links section in the Slack messages. Default: true.

* `show_section_alert_details` - (Optional) Whether to include the alert details section in the Slack messages. Default: true.

* `show_section_capturing_information` - (Optional) Whether to include the capturing information section in the Slack messages. Default: true.

* `channel` - (Required) Channel name from this Slack.

* `is_private_channel` - (Optional) Whether the Slack channel is private or not.

* `private_channel_url` - (Optional) The channel URL, if channel is private.


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

* `share_with_current_team` - (Optional) If set to `true` it will share notification channel only with current team (in which user is logged in).
  Otherwise, it will share it with all teams, which is the default behaviour. Although this is an optional setting, beware that if you have lower permissions than admin you may see a `error: 403 Forbidden` if this is not set to `true`.

## Import

Slack notification channels for Monitor can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_notification_channel_slack.example 12345
```
