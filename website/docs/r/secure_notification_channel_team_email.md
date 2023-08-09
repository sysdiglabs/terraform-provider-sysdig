---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel_team_email"
description: |-
  Creates a Sysdig Secure Notification Channel of type Team Email.
---

# Resource: sysdig_secure_notification_channel_team_email

Creates a Sysdig Secure Notification Channel of type Team Email.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_notification_channel_team_email" "sample-gchat" {
  name                    = "Example Channel - team email"
  team_id                 = 1
  enabled                 = true
  notify_when_ok          = false
  notify_when_resolved    = false
  share_with_current_team = true
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `team_id` - (Required) id of the team.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.

* `share_with_current_team` - (Optional) If set to `true` it will share notification channel only with current team (in which user is logged in).
  Otherwise, it will share it with all teams, which is the default behaviour.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the Notification Channel.

* `version` - (Computed) The current version of the Notification Channel.

## Import

Team email notification channels for Secure can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_notification_channel_team_email.example 12345
```
