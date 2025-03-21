---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_notification_channel_ibm_function"
description: |-
  Creates a Sysdig Monitor Notification Channel of type IBM Function.
---

# Resource: sysdig_monitor_notification_channel_ibm_function

Creates a Sysdig Monitor Notification Channel of type IBM Function.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_notification_channel_ibm_function" "sample-ibm-function-web-action" {
  name              = "Example Channel - IBM Function - web action"
  enabled           = true
  ibm_function_type = "WEB_ACTION"
  url               = "https://eu-gb.functions.cloud.ibm.com/api/v1/web/namespaces/eeeeeeee-623b-4776-ba35-4065bcbfee7b/actions/hello-world/helloworld?param=true"
  whisk_auth_token = "xxx"

  custom_data = {
    "data1": "value1"
    "data2": "value2"
  }

  notify_when_ok          = false
  notify_when_resolved    = false
  send_test_notification  = false
}

resource "sysdig_monitor_notification_channel_ibm_function" "sample-ibm-function-cloud-function" {
  name              = "Example Channel - IBM Function - cloud function"
  ibm_function_type = "CLOUD_FUNCTION"
	url               = "https://eu-gb.functions.cloud.ibm.com/api/v1/namespaces/13eeeeee-623b-4776-ba35-4065bcbfee7b/actions/hello-world/myaction"
	iam_api_key       = "xxx"
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `ibm_function_type` - (Required) Type of IBM Function. Can be `WEB_ACTION` for a Web Action (with or without X-Require-Whisk-Auth header) or `CLOUD_FUNCTION` for an IAM Secured Action.

* `url` - (Required) URL of the IBM Function.

* `custom_data` - (Optional) Key value list of additional parameters for the IBM Function.

* `whisk_auth_token` - (Optional) Only if `ibm_function_type` is `WEB_ACTION`: Whisk authentication token.

* `iam_api_key` - (Optional) Required if `ibm_function_type` is `CLOUD_FUNCTION`: API Key to call the private cloud function.

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

IBM Function notification channels for Monitor can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_notification_channel_ibm_function.example 12345
```
