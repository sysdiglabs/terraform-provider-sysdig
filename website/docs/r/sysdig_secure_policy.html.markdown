---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_policy"
sidebar_current: "docs-sysdig-secure-policy"
description: |-
  Creates a Sysdig Secure Policy.
---

# sysdig\secure\_policy

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
  //filter = "host.ip.private = \"10.0.23.1\""
  container_scope = true
  host_scope = true

  //actions {
  //  container = "pause"

  //  capture {
  //    seconds_before_event = 60
  //    seconds_after_event = 60
  //  }
  //}

  // Falco rule selection
  falco_rule_name_regex = "Unexpected spawned process traefik"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure policy. It must be unique.

* `description` - (Required) The description of Secure policy.

* `severity` - (Required) The severity of Secure policy. The accepted values
    are: 2 (High), 4 (Medium) and 6 (Low).

* `enabled` - (Required) Will secure process with this rule?

- - -

### Scope selection

* `host_scope` - (Required) The application scope of this rule. Does this rule
    applies to hosts?

* `container_scope` - (Required) The application scope of this rule. Does this
    rule applies to containers? Note that the rule should apply at least to one
    scope, host or container.

* `filter` - (Optional) Limit appplication scope based in one expresion. By
    example: "host.ip.private = \"10.0.23.1\""

- - -

### Actions block

The actions block is optional and supports:

* `container` - (Required) The action applied to container when this Policy is
    triggered. Can be *stop* or *pause*.

which
The capture block is optional and whan present captures with Sysdig the stream
of system calls:

* `seconds_before_event` - (Required) Captures the system calls during the
    amount of seconds before the policy was triggered.

* `seconds_after_event` - (Required) Captures the system calls for the amount
    of seconds after the policy was triggered.

- - -

### Falco rule selection

* `falco_rule_name_regex` - (Required) The RegExp for checking matches with
    Falco Rule name.  When a you have uploaded custom rules, and an alert is
    raised. Check if that alert matches with this regexp for raising the Policy
    Alert.
