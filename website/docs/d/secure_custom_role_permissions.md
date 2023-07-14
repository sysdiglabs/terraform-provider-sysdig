---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_custom_role_permissions"
description: |-
  Validate and enrich with permissions on which the requested permissions depend
---

# Data Source: sysdig_secure_custom_role_permissions

Validate and enrich the requested permissions


-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_custom_role_permissions" "images_edit" {
  requested_permissions = ["secure.blacklist.images.edit"]
}
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `enriched_permissions` - The minimum set of permissions enriched with permissions on which the requested permissions depend


## Example Usage with Custom Role

```terraform
data "sysdig_secure_custom_role_permissions" "images_edit" {
  requested_permissions = ["secure.blacklist.images.edit"]
}

resource "sysdig_custom_role" "my-custom-role" {
  depends_on = [data.sysdig_secure_custom_role_permissions.images_edit]
  name = "custom-role-name"
  description = "Custom role to edit images"

  permissions {
    secure_permissions = data.sysdig_secure_custom_role_permissions.images_edit.enriched_permissions
  }
}
```
