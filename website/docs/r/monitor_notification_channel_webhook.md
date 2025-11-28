---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_webhook"
description: |-
  Creates a Sysdig Monitor Notification Channel of type Webhook.
---

# Resource: sysdig_monitor_notification_channel_webhook

Creates a Sysdig Monitor Notification Channel of type Webhook.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_notification_channel_webhook" "sample-webhook" {
  name                    = "Example Channel - Webhook"
  enabled                 = true
  url                     = "localhost:8080"
  send_test_notification  = false

  custom_data = {
    "data1": "value1"
    "data2": "value2"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `url` - (Required) URL to send the event.

* `custom_data` - (Optional) Key value list of additional data you want to attach to the alert notification.

* `additional_headers` - (Optional) Key value list of custom headers.

* `allow_insecure_connections` - (Optional) Whether to skip TLS verification. Default: `false`.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional, Deprecated) Send a new notification when the alert condition is no longer triggered. Default is `true`. This option is deprecated; use `notify_on_resolve` within the `notification_channels` options in the `sysdig_monitor_alert_v2_*` resources instead, which takes precedence over this setting.

* `notify_when_resolved` - (Optional, Deprecated) Send a new notification when the alert is manually acknowledged by a user. Default is `true`. This option is deprecated; use `notify_on_acknowledge` within the `notification_channels` options in the `sysdig_monitor_alert_v2_*` resources instead, which takes precedence over this setting.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the Notification Channel.

* `version` - (Computed) The current version of the Notification Channel.

* `share_with_current_team` - (Optional) If set to `true` it will share notification channel only with current team (in which user is logged in).
  Otherwise, it will share it with all teams, which is the default behaviour. Although this is an optional setting, beware that if you have lower permissions than admin you may see a `error: 403 Forbidden` if this is not set to `true`.

## Import

Webhook notification channels for Monitor can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_notification_channel_webhook.example 12345
```
