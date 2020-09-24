---
layout: "sysdig"
page_title: "Sysdig: sysdig_current_user"
sidebar_current: "docs-sysdig-current-user-ds"
description: |-
  Retrieves information about the user performing the API calls.
---

# sysdig\_current\_user

Retrieves information about the current user performing the API calls.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
data "sysdig_current_user" "me" {
}
```

## Attributes Reference

* `id` - The current user's ID.

* `email` - The user email.

* `name` - The user's first name.

* `last_name` - The user's last name.

* `system_role` - The user's system role.