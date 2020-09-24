---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_team"
sidebar_current: "docs-sysdig-secure-team"
description: |-
  Creates a Sysdig Secure Team.
---

# sysdig\_secure\_team

Creates a Sysdig Secure Team.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_team" "devops" {
  name = "DevOps team"
  
  user_roles {
    email = data.sysdig_current_user.me.email
    role = "ROLE_TEAM_MANAGER"
  }

  user_roles {
    email = "john.doe@example.com"
    role = "ROLE_TEAM_STANDARD"
  }
}
 
data "sysdig_current_user" "me" {
}
```

## Argument Reference

* `name` - (Required) The name of the Secure Team. It must be unique.

* `description` - (Optional) A description of the team.

* `theme` - (Optional) Colour of the team. Default: "#73A1F7".

* `scope_by` - (Optional) Scope for the team. Default: "container".

* `filter` - (Optional) If the team can only see some resources, 
             write down a filter of such resources.
             
* `use_sysdig_capture` - (Optional) Defines if the team is able to create Sysdig Capture files. 
                         Default: true.
                         
* `default_team` - (Optional) Defines if the team is the default one. Warning: only one can be the default,
                   if you define multiple default teams, Terraform will be updating the API in every execution,
                   even if the state hasn't changed.

* `user_roles` - (Optional) Multiple user roles can be specified.
                 Administrators of the account will be automatically added
                 to every new created team, so they don't need to be added as a
                 resource in the Terraform manifest.
                         
### User Role Argument Reference

* `email` - (Required) The email of the user in the group.

* `role` - (Optional) The role for the user in this group.
           Valid roles are: ROLE_TEAM_STANDARD, ROLE_TEAM_EDIT, ROLE_TEAM_READ, ROLE_TEAM_MANAGER.
           Default: ROLE_TEAM_STANDARD.