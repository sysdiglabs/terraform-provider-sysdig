---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_group_mapping"
description: |-
  Creates a group mapping in Sysdig.
---

# Resource: sysdig_group_mapping

Creates a group mapping in Sysdig.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

### Regular users

```terraform
resource "sysdig_group_mapping" "my_group" {
  group_name = "my-group"
  role = "ROLE_TEAM_STANDARD"
  system_role = "ROLE_USER"

  team_map {
    all_teams = false
    team_ids = [sysdig_secure_team.my_team.id, sysdig_monitor_team.my_team.id]
  }
  weight = 10
}

```
This way, we define a group mapping named "my-group" for a user who will have a standard role in two teams.

### Admin users
If the group members should assume the Sysdig administrator role the mapping should be created this way

```terraform
resource "sysdig_group_mapping" "admin" {
  group_name = "admin"
  role = "ROLE_TEAM_MANAGER"
  system_role = "ROLE_CUSTOMER"

  team_map {
    all_teams = true
    team_ids = []
  }
}
```
The name doesn’t necessarily have to be “admin,” it’s just an example. The important aspects are the roles and the team_map

## Argument Reference

* `group_name` - (Required) The group name to be mapped.

* `role` - (Required) The role that is assigned to the users. It can be a standard role or a custom team role ID.

* `system_role` (Optional) The system role that is assigned to the users. The supported values are: 
  * `ROLE_USER` for regular users (Default if not specified) 
  * `ROLE_CUSTOMER` for admin users

* `team_map` - (Required) Block to define team mapping.

* `weight` - (Optional) The group mapping weight used to solve conflicts. Weight is a positive number, lower number has higher priority.

### team_map

Team map block is required and supports:

* `all_teams` - (Optional) Flag indicating whether team map should resemble all customer teams.

* `team_ids` - (Optional) Set of team IDs, is empty when `all_teams` is true, otherwise needs at least 1 element.


## Attributes Reference

No additional attributes are exported.

## Import

Sysdig group mapping can be imported using the ID, e.g.

```
$ terraform import sysdig_group_mapping.my_group 24267
```
