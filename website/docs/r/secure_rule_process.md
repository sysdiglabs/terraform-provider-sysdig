---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_process"
description: |-
  Creates a Sysdig Secure Process Rule.
---

# Resource: sysdig_secure_rule_process

Creates a Sysdig Secure Process Rule.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_secure_rule_process" "sample" {
  name = "Launch Suspicious Network Tool in Container" // ID
  description = "Detect network tools launched inside container"

  matching = true // default
  processes = ["nc", "ncat", "nmap", "dig", "tcpdump", "tshark", "ngrep"]
}

```

## Argument Reference

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Optional) The description of Secure rule. By default is empty.
* `tags` - (Optional) A list of tags for this rule.

### Matching

* `matching` - (Optional) Defines if the process name matches or not with the provided list. Default is true.
* `processes` - (Required) List of processes to match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

## Import

Secure process runtime rules can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_rule_process.example 12345
```