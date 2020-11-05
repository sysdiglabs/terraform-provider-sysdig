---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_filesystem"
sidebar_current: "docs-sysdig-secure-rule-filesystem"
description: |-
  Creates a Sysdig Secure Filesystem Rule.
---

# sysdig\_secure\_rule\_filesystem

Creates a Sysdig Secure Filesystem Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl

resource "sysdig_secure_rule_filesystem"  "example" {
  name = "Apache writing to non allowed directory"
  description = "Attempt to write to directories that should be immutable"
  tags = ["filesystem", "cis"]

  read_only {
    matching = true // default
    paths = ["/etc"]
  }

  read_write {
    matching = true // default
    paths = ["/var/log/apache2", "/dev/tty"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Optional) The description of Secure rule. By default is empty.
* `tags` - (Optional) A list of tags for this rule.

### Read Only

* `matching` - (Optional) Defines if the path matches or not with the provided list. Default is true.
* `paths` - (Required) List of paths to match.

### Read Write

* `matching` - (Optional) Defines if the path matches or not with the provided list. Default is true.
* `paths` - (Required) List of paths to match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.