---
layout: "sysdig"
page_title: "Sysdig: sysdig_current_user"
description: |-
  Retrieves information about the user performing the API calls.
---

# Data Source: sysdig\_current\_user

Retrieves information about the current user performing the API calls.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

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