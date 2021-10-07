---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel_victorops"
description: |-
  Creates a Sysdig Secure Notification Channel of type VictorOps.
---

# Resource: sysdig_secure_notification_channel_victorops

Creates a Sysdig Secure Notification Channel of type VictorOps.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example Usage

```hcl
resource "sysdig_secure_notification_channel_victorops" "sample-victorops" {
	name                    = "Example Channel - VictorOps"
	enabled                 = true
	api_key                 = "1234342-4234243-4234-2"
	routing_key             = "My team"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `api_key` - (Required) Key for the API.

* `routing_key` - (Required) Routing key for VictorOps. 

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

VictorOPS notification channels for Secure can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_notification_channel_victorops.example 12345
```