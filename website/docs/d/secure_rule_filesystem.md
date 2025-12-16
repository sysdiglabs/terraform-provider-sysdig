---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_filesystem"
description: |-
  Retrieves a Sysdig Secure Filesystem Rule.
---

# Data Source: sysdig_secure_rule_filesystem

Retrieves the information of an existing Sysdig Secure Filesystem Rule.

~> **DEPRECATED:** List matching rules have been deprecated. Please use [sysdig_secure_rule_falco](../d/secure_rule_falco.html) instead for reading runtime security rules.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_rule_file_system" "example" {
    name = "Write below etc"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure rule to retrieve.

## Attributes Reference

In addition to the argument above, the following attributes are exported:

* `description` - The description of Secure rule.
* `tags` - A list of tags for this rule.
* `read_only` - Block that defines read only paths to match or not match.
* `read_write` - Block that defines read and write paths to match or not match.
* `version` - Current version of the resource in Sysdig Secure.

## read_write and read_only blocks

Description of the attributes within the read_only and read_write blocks.

* `matching` - Boolean value that defines if the path matches or not with the provided list.
* `paths` - List of paths to match.