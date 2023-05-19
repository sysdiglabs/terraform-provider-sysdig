---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_process"
description: |-
  Retrieves a Sysdig Secure Process Rule.
---

# Data Source: sysdig_secure_rule_process

Retrieves the information of an existing Sysdig Secure Process Rule.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_rule_process" "example" {
    name = "Launch Suspicious Network Tool in Container"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure rule to retrieve.

## Attributes Reference

In addition to the argument above, the following attributes are exported:

* `description` - The description of Secure rule.
* `tags` - A list of tags for this rule.
* `matching` - Defines if the process name matches or not with the provided list.
* `processes` - List of processes to match.
* `version` - Current version of the resource in Sysdig Secure.
