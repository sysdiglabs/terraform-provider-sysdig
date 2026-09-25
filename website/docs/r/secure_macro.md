---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_macro"
description: |-
  Creates a Sysdig Secure Falco Macro.
---

# Resource: sysdig_secure_macro

Creates a Sysdig Secure Falco Macro.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_macro" "web_port" {
  name      = "web_port"
  condition = "fd.sport=80"
}
```

To extend a macro that Sysdig ships, create a macro with the same `name` and
`append = true`. The condition is concatenated onto the existing one, so it must
begin with a logical operator (`or` / `and`):

```terraform
resource "sysdig_secure_macro" "allow_my_etc_writer" {
  name      = "user_known_write_below_etc_activities"
  condition = "or (proc.name = my_installer and fd.name startswith /etc/myapp/)"
  append    = true # default: false
}
```

## Argument Reference

* `name` - (Required) The name of the macro. It must be unique if it's not in append mode.

* `condition` - (Required) Macro condition. It can contain lists or other macros.

* `append` - (Optional) Adds these elements to an existing macro. Used to extend existing macros provided by Sysdig.
    A macro of the same `name` must already exist, and `condition` must begin with a logical operator (`or` / `and`)
    because it is concatenated onto the existing condition. By default this is false.

    ~> **Note:** Appending a macro that you created yourself (rather than one provided by Sysdig) is only supported on
    backends that have migrated to the v2 macro storage. On earlier backends the create fails with
    `The field 'name' must not be the same as another Secure UI macro`. Appending a Sysdig-provided macro works on all
    backends. On earlier backends a macro can also only be extended once; newer backends do not enforce that limit.

* `minimum_engine_version` - (Optional) This is used to indicate that the macro requires a minimum engine version. This
    can allow you to add macros that would not normally pass validation with older agents in your environment. The macro
    will only be processed by agents that support the minimum_engine_version specified.


## Attributes Reference

No additional attributes are exported.

## Import

Secure macros can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_macro.example 12345
```