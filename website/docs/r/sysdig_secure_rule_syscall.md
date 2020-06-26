---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_syscall"
sidebar_current: "docs-sysdig-secure-rule-syscall"
description: |-
  Creates a Sysdig Secure Syscall Rule.
---

# sysdig\_secure\_rule\_syscall

Creates a Sysdig Secure Syscall Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_rule_syscall" "foo" {
  name = "Unexpected mount syscall" // ID
  description = "Syscall 'mount' detected"

  matching = true // default
  syscalls = ["mount"]
}
```

## Argument Reference

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Required) The description of Secure rule.
* `tags` - (Optional) A list of tags for this rule.

### Matching

* `matching` - (Optional) Defines if the syscall name matches or not with the provided list. Default is true.
* `processes` - (Required) List of syscalls to match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.