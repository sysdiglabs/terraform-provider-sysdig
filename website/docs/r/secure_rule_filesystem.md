---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_filesystem"
description: |-
  Creates a Sysdig Secure Filesystem Rule.
---

# Resource: sysdig_secure_rule_filesystem

Creates a Sysdig Secure Filesystem Rule.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example Usage

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

## Import

Secure filesystem runtime rules can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_rule_filesystem.example 12345
```