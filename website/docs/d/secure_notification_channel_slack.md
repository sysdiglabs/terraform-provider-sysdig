---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel_slack"
description: |-
  Retrieves information about a Secure notification channel of type Slack
---

# Data Source: sysdig_secure_notification_channel_slack

Retrieves information about a Secure notification channel of type Slack.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_notification_channel_slack" "nc_slack" {
	name = "some notification channel name"
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Notification Channel ID.
* `url` - URL of the Slack webhook.
* `channel` - Name of the Slack channel.
* `is_private_channel` - Whether the Slack Channel has been marked as private or not.
* `private_channel_url` - The channel URL, i.e. the link that is referencing the channel (not to be confused with the webhook url), if the channel is private.
* `template_version` - The notification template version to use to create notifications.
* `enabled` - Whether the Notification Channel is active or not.
* `notify_when_ok` - Whether the Notification Channel sends a notification when the condition is no longer triggered.
* `notify_when_resolved` - Whether the Notification Channel sends a notification if it's manually acknowledged by a user.
* `version` - The version of the Notification Channel.
* `send_test_notification` - Whether the Notification Channel has enabled the test notification.
