---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_managed_policy"
description: |-
  Manages configuration of a Sysdig Secure Managed Policy.
---

# Resource: sysdig_secure_managed_policy

Manages configuration of a Sysdig Secure Managed Policy.

-> **Note:** Sysdig managed policies are not resources that you create. They are provided by Sysdig. This resource
allows you to identify and configure a managed policy. The managed policy is looked up by its name and type.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_managed_policy" "sysdig_runtime_threat_detection" {
  name = "Sysdig Runtime Threat Detection"
  type = "falco"
  enabled = true
  runbook = "https://runbook.com"
  
  // Scope selection
  scope = "container.id != \"\""

  // Disabling rules
  disabled_rules = ["Suspicious Cron Modification"]

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
    }
  }

  notification_channels = [10000]
}
```

## Argument Reference

* `name` - (Required) The name of the Secure managed policy. It must match the name of an existing managed policy.

* `type` - (Optional) Specifies the type of the runtime policy. Must be one of: `falco`, `list_matching`, `k8s_audit`,
  `aws_cloudtrail`, `gcp_auditlog`, `azure_platformlogs`. By default it is `falco`.

* `enabled` - (Optional) Will secure process with this policy?. By default this is true.

* `runbook` - (Optional) Customer provided url that provides a runbook for a given policy. 

- - -

### Scope selection

* `scope` - (Optional) Limit application scope based in one expression. For
    example: "host.ip.private = \\"10.0.23.1\\"". By default the rule won't be scoped
    and will target the entire infrastructure.

- - -

### Actions block

The actions block is optional and supports:

* `container` - (Optional) The action applied to container when this Policy is
    triggered. Can be *stop*, *pause* or *kill*. If this is not specified,
    no action will be applied at the container level.

* `capture` - (Optional) Captures with Sysdig the stream of system calls:
    * `seconds_before_event` - (Required) Captures the system calls during the
    amount of seconds before the policy was triggered.
    * `seconds_after_event` - (Required) Captures the system calls for the amount
    of seconds after the policy was triggered.
    * `name` - (Optional) The name of the capture file

- - -

### Disabling falco rules

* `disabled_rules` - (Optional) Array with the name of the rules in the managed policy to disable.

- - -

### Notification

* `notification_channels` - (Optional) IDs of the notification channels to send alerts to
    when the policy is fired.

## Attributes Reference

No additional attributes are exported.
