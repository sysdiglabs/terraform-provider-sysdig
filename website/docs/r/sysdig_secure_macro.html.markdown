---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_macro"
sidebar_current: "docs-sysdig-secure-macro"
description: |-
  Creates a Sysdig Secure Falco Macro.
---

# sysdig\_secure\_macro

Creates a Sysdig Secure Falco Macro.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

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
