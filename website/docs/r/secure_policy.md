---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_policy"
description: |-
  Creates a Sysdig Secure Policy.
---

# Resource: sysdig_secure_policy

Creates a Sysdig Secure Policy.

~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.  

## Example Usage

```hcl
resource "sysdig_secure_policy" "write_apt_database" {
  name = "Write apt database"
  description = "an attempt to write to the dpkg database by any non-dpkg related program"
  severity = 4
  enabled = true

  // Scope selection
  scope = "container.id != \"\""

  // Rule selection
  rule_names = ["Terminal shell in container"]

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

* `name` - (Required) The name of the Secure policy. It must be unique.

* `description` - (Required) The description of Secure policy.

* `severity` - (Optional) The severity of Secure policy. The accepted values
    are: 0, 1, 2, 3 (High), 4, 5 (Medium), 6 (Low) and 7 (Info). The default value is 4 (Medium).

* `enabled` - (Optional) Will secure process with this rule?. By default this is true.

* `type` - (Optional) Specifies the type of the runtime policy. Must be one of: `falco`, `list_matching`, `k8s_audit`, `aws_cloudtrail`. By default it is `falco`.

- - -

### Scope selection

* `scope` - (Optional) Limit appplication scope based in one expresion. For
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

- - -

### Falco rule selection

* `rule_names` - (Optional) Array with the name of the rules to match.

- - -

### Notification

* `notification_channels` - (Optional) IDs of the notification channels to send alerts to
    when the policy is fired.

## Import

Secure runtime policies can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_policy.example 12345
```