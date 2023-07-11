---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_custom_role"
description: |-
  Retrieves information about a custom role from the name
---

# Data Source: sysdig_custom_role

Retrieves information about a custom role from the name

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_custom_role" "custom_role" {
  name = "CustomRoleName"
}
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The custom role's ID.

* `name` - The custom role's name.

* `description` - The custom role's description.

* `monitor_permissions` - The custom role's monitor permissions.

* `secure_permissions` - The custom role's secure permissions.
