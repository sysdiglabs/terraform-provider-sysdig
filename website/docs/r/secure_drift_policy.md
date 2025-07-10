---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_drift_policy"
description: |-
  Retrieves a Sysdig Secure Drift Policy.
---

# Resource: sysdig_secure_drift_policy

Retrieves the information of an existing Sysdig Secure Drift Policy.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_notification_channel" "email_notification_channel" {
  name = "Test Email Channel"
}

resource "sysdig_secure_drift_policy" "policy" {
  name = "Drift Policy 1"
  description = "<Description>"
  severity = 4
  enabled = true
  runbook = "https://runbook.com"
  
  // Scope selection
  scope = "container.id != \"\""

  // Rule selection
  rule {
    description = "Test Drift Rule Description"
    enabled     = true

    exceptions {
      items       = ["/usr/bin/curl"]
      match_items = false
    }

    prohibited_binaries {
      items       = ["/usr/bin/sh"]
      match_items = true
    }

  }

  actions {
    prevent_drift = true
    container     = "stop"

    capture {
      seconds_before_event = 5
      seconds_after_event  = 10
    }
  }

  notification_channels = [data.sysdig_secure_notification_channel.email_notification_channel.id]
}
```

## Argument Reference

* `name` - (Required) The name of the Secure managed policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id for the policy.

* `description` - The description for the managed policy.

* `severity` -  The severity of Secure policy. The accepted values
    are: 0, 1, 2, 3 (High), 4, 5 (Medium), 6 (Low) and 7 (Info).

* `enabled` - Whether the policy is enabled or not.

* `runbook` - Customer provided url that provides a runbook for a given policy.

* `scope` - The application scope for the policy.

* `notification_channels` - IDs of the notification channels to send alerts to
    when the policy is fired.

### Actions block

The actions block is optional and supports the following:

* `prevent_drift` - (Optional) Prevent the execution of drifted binaries and specified prohibited binaries.

For agents 12.20 and above, these additional actions are supported: 

* `container` - (Optional) The action applied to container when this Policy is
    triggered. Can be *stop*, *pause* or *kill*. If this is not specified,
    no action will be applied at the container level.

* `capture` - (Optional) Captures with Sysdig the stream of system calls:
    * `seconds_before_event` - (Required) Captures the system calls during the
    amount of seconds before the policy was triggered.
    * `seconds_after_event` - (Required) Captures the system calls for the amount
    of seconds after the policy was triggered.
    * `name` - (Required) The name of the capture file
    * `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name=cat`
    * `bucket_name` - (Optional) Custom bucket to store capture in, 
    bucket should be onboarded in Integrations > S3 Capture Storage. Default is to use Sysdig Secure Storage 
    * `folder` - (Optional) Name of folder to store capture inside the bucket. 
    By default we will store the capture file at the root of the bucket

### `rule` block

The rule block is required and supports:

* `description` - (Required) The description of the drift rule.
* `enabled` - (Required) Toggle to dynamically detect execution of drifted binaries. A drifted binary is any binary that was not part of the original image of the container. It is typically downloaded or compiled into a running container.
* `exceptions` - (Optional) Specify comma separated list of exceptions.
    * `items` - (Required) Specify comma separated list of exceptions, e.g. `/usr/bin/rm, /usr/bin/curl`.
* `prohibited_binaries` - (Optional) A prohibited binary can be a known harmful binary or one that facilitates discovery of your environment.
    * `items` - (Required) Specify comma separated list of prohibited binaries, e.g. `/usr/bin/rm, /usr/bin/curl`.
* `process_based_exceptions` - (Optional) List of processes that will be able to execute a drifted file 
    * `items` - (Required) Specify comma separated list of processes, e.g. `/usr/bin/rm, /usr/bin/curl`.      
* `process_based_prohibited_binaries` - (Optional) List of processes that will be prohibited to execute a drifted file
    * `items` - (Required) Specify comma separated list of processes, e.g. `/usr/bin/rm, /usr/bin/curl`.
* `mounted_volume_drift_enabled` - (Optional) Treat all binaries from mounted volumes as drifted. Default value is false/disabled.
* `use_regex` - (Optional) Pass exceptions and prohibited binaries as regex strings. Requires agent version 13.2.0 and above
