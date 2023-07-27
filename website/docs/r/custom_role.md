---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_custom_role"
description: |-
  Creates a custom role in Sysdig.
---

# Resource: sysdig_custom_role

Creates a custom role in Sysdig.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_custom_role" "my-custom-role" {
  name = "custom-role-name"
  description = "test"

  permissions {
    monitor_permissions = ["kubernetes-api-commands.read"]
    secure_permissions = ["scanning.read"]
  }
}

```

## Argument Reference

* `name` - (Required) The custom role name.

* `description` - (Optional) Additional long description.

* `permissions` (Required) Block to define monitor and secure permissions.

### permissions

Permissions block is required and supports:

* `monitor_permissions` - (Optional) Set of Monitor permissions assigned to the role. Check GET /api/permissions to get the list of available values

* `secure_permissions` - (Optional) Set of Secure permissions assigned to the role. Check GET /api/permissions to get the list of available values.

### Permissions data source

Permissions can have dependencies and dependee. Since the dependencies graph can be hard to determine manually were introduced 

[`sysdig_secure_custom_role_permissions`](../d/secure_custom_role_permissions.md) and [`sysdig_monitor_custom_role_permissions`](../d/monitor_custom_role_permissions.md)

Please check the relative documentation to see how to use them.

## Attributes Reference

No additional attributes are exported.

## Import

Sysdig group mapping can be imported using the ID, e.g.

```
$ terraform import sysdig_custom_role.my_custom_role 50
```
