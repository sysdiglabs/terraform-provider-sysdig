---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_builtin_role"
description: |-
  Retrieves information about a built-in (OOTB) role from the name.
---

# Data Source: sysdig_builtin_role

Retrieves information about a built-in (out-of-the-box) role from the name.

Built-in roles are the roles provided by Sysdig: View Only, Standard User, Advanced User, and Team Manager.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_builtin_role" "advanced_user" {
  name = "Advanced User"
}
```

## Argument Reference

* `name` - (Required) The name of the built-in role. Valid values are: `View Only`, `Standard User`, `Advanced User`, `Team Manager`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `monitor_permissions` - The built-in role's monitor permissions.

* `secure_permissions` - The built-in role's secure permissions.
