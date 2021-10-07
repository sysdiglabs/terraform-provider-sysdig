---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_user"
description: |-
  Creates a user in Sysdig.
---

# Resource: sysdig_user

Creates a user in Sysdig.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example Usage

```terraform
resource "sysdig_user" "foo-bar" {
  email = "foo.bar@sysdig.com"
  system_role = "ROLE_CUSTOMER"
  first_name = "foo"
  last_name = "bar"
}
```

## Argument Reference

* `email` - (Required) The email for the user to invite.

* `system_role` - (Optional) The privileges for the user. It can be either "ROLE_USER" or "ROLE_CUSTOMER".
    If set to "ROLE_CUSTOMER", the user will be known as an admin.

* `first_name` - (Optional) The name of the user.

* `last_name` - (Optional) The last name of the user.


## Attributes Reference

No additional attributes are exported.

## Import

Sysdig users can be imported using the ID, e.g.

```
$ terraform import sysdig_user.example 12345
```