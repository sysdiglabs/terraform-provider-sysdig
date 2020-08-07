---
layout: "sysdig"
page_title: "Sysdig: sysdig_user"
sidebar_current: "docs-sysdig-user"
description: |-
  Creates a user in Sysdig.
---

# sysdig\_user

Creates a user in Sysdig.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
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

* `password` - (Optional) The password for the user. If the password is defined, the user will not receive
  a confirmation email.
