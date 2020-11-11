---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_policy"
sidebar_current: "docs-sysdig-secure-policy"
description: |-
  Creates a Sysdig Secure Policy.
---

# sysdig\_secure\_policy

Creates a Sysdig Secure Policy.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

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
