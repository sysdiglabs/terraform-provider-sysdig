---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_default_role"
description: |-
  Retrieves information about a default (OOTB) role from the name.
---

# Data Source: sysdig_default_role

Retrieves information about a default (out-of-the-box) role from the name.

Default roles are the built-in roles provided by Sysdig: View Only, Standard User, Advanced User, and Team Manager.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_default_role" "advanced_user" {
  name = "Advanced User"
}
```

## Argument Reference

* `name` - (Required) The name of the default role. Valid values are: `View Only`, `Standard User`, `Advanced User`, `Team Manager`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `monitor_permissions` - The default role's monitor permissions.

* `secure_permissions` - The default role's secure permissions.
