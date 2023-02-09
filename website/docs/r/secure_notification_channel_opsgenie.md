---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel"
description: |-
  Creates a Sysdig Secure Notification Channel of type OpsGenie.
---

# Resource: sysdig_secure_notification_channel_opsgenie

Creates a Sysdig Secure Notification Channel of type OpsGenie.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_notification_channel_opsgenie" "sample-opsgenie" {
	name                    = "Example Channel - OpsGenie"
	enabled                 = true
	api_key                 = "2349324-342354353-5324-23"
	notify_when_ok          = false
	notify_when_resolved    = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `api_key` - (Required) Key for the API.

* `region` - (Optional) Opsgenie Region. Can be `US` or `EU`. Default is `US`.

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

Opsgenie notification channels for Secure can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_notification_channel_opsgenie.example 12345
```
