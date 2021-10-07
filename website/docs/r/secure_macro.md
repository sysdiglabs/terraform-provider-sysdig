---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_macro"
description: |-
  Creates a Sysdig Secure Falco Macro.
---

# Resource: sysdig_secure_macro

Creates a Sysdig Secure Falco Macro.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example Usage

```hcl
resource "sysdig_secure_macro" "http_port" {
  name = "web_port"
  condition = "fd.sport=80"
}

resource "sysdig_secure_macro" "https_port" {
  name = "web_port"
  condition = "or fd.sport=443"
  append = true # default: false
}
```

## Argument Reference

* `name` - (Required) The name of the macro. It must be unique if it's not in append mode.

* `condition` - (Required) Macro condition. It can contain lists or other macros.

* `append` - (Optional)  Adds these elements to an existing macro. Used to extend existing macros provided by Sysdig.
    The macros can only be extended once, for example if there is an existing macro called "foo", one can have another 
    append macro called "foo" but not a second one. By default this is false.

## Import

Secure macros can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_macro.example 12345
```