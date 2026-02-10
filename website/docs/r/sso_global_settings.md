---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_sso_global_settings"
description: |-
  Manages global SSO settings per product in Sysdig using the Platform API.
---

# Resource: sysdig_sso_global_settings

Manages global SSO settings for a specific Sysdig product (Monitor or Secure) using the Platform API.

This is a singleton resource per product — only one instance should exist for each product. The resource cannot be deleted; removing it from Terraform configuration will only remove it from state.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_sso_global_settings" "monitor" {
  product                   = "monitor"
  is_password_login_enabled = true
}

resource "sysdig_sso_global_settings" "secure" {
  product                   = "secure"
  is_password_login_enabled = false
}
```

## Argument Reference

* `product` - (Required, ForceNew) The Sysdig product. Valid values: `monitor`, `secure`. Changing this forces creation of a new resource.

* `is_password_login_enabled` - (Required) Whether password-based login is enabled alongside SSO for this product.

## Attributes Reference

No additional attributes are exported.

## Import

SSO global settings can be imported using the product name:

```
$ terraform import sysdig_sso_global_settings.monitor monitor
$ terraform import sysdig_sso_global_settings.secure secure
```
