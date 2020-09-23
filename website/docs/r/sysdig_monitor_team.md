---
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_team"
sidebar_current: "docs-sysdig-monitor-team"
description: |-
  Creates a Sysdig Monitor Team.
---

# sysdig\_monitor\_team

Creates a Sysdig Monitor Team.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_monitor_team" "devops" {
  name = "Monitoring DevOps team"

  entrypoint {
	type = "Explore"
  }
  
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

* `name` - (Required) The name of the Monitor Team. It must be unique.

* `user_roles` - (Required) At least one user role must be specified, which in most of the
                 scenarios will be the calling user, as well as other accounts.
                 
* `entrypoint` - (Required) Main entry point for the current team in the product. 
                 See the Entrypoint argument reference section for more information.

* `description` - (Optional) A description of the team.

* `theme` - (Optional) Colour of the team. Default: "#73A1F7".

* `scope_by` - (Optional) Scope for the team. Default: "container".

* `filter` - (Optional) If the team can only see some resources, 
             write down a filter of such resources.
             
* `use_sysdig_capture` - (Optional) Defines if the team is able to create Sysdig Capture files. 
                         Default: true.
                         
* `can_see_infrastructure_events` - (Optional) TODO. Default: false.

* `can_use_aws_data` - (Optional) TODO. Default: false.

### Entrypoint Argument Reference

* `type` - (Required) Main entrypoint for the team.
                      Valid options are: Explore, Dashboards, Events, Alerts, Settings.

* `selection` - (Optional) Sets up the defined Dashboard name as entrypoint.
                Warning: This field must only be added if the `type` is "Dashboards".

### User Role Argument Reference

* `email` - (Required) The email of the user in the group.

* `role` - (Optional) The role for the user in this group.
           Valid roles are: ROLE_TEAM_STANDARD, ROLE_TEAM_EDIT, ROLE_TEAM_READ, ROLE_TEAM_MANAGER.
           Default: ROLE_TEAM_STANDARD.
                 