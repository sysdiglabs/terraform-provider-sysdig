---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_sso_group_mapping"
description: |-
  Creates an SSO group mapping in Sysdig using the Platform API.
---

# Resource: sysdig_sso_group_mapping

Creates an SSO group mapping in Sysdig using the Platform API. This resource replaces the deprecated `sysdig_group_mapping` resource.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

### Standard role for all teams

```terraform
resource "sysdig_sso_group_mapping" "all_teams" {
  group_name         = "engineering"
  standard_team_role = "ROLE_TEAM_STANDARD"
  is_admin           = false

  team_map {
    is_for_all_teams = true
  }

  weight = 10
}
```

### Custom role for specific teams

```terraform
resource "sysdig_sso_group_mapping" "specific_teams" {
  group_name          = "devops"
  custom_team_role_id = sysdig_custom_role.devops_role.id

  team_map {
    is_for_all_teams = false
    team_ids         = [sysdig_secure_team.my_team.id, sysdig_monitor_team.my_team.id]
  }

  weight = 20
}
```

### Admin group mapping

```terraform
resource "sysdig_sso_group_mapping" "admins" {
  group_name         = "platform-admins"
  standard_team_role = "ROLE_TEAM_MANAGER"
  is_admin           = true

  team_map {
    is_for_all_teams = true
  }
}
```

## Argument Reference

* `group_name` - (Required) The SSO group name to map. Maximum 256 characters.

* `standard_team_role` - (Optional) The standard team role assigned to users. Conflicts with `custom_team_role_id`. One of `standard_team_role` or `custom_team_role_id` must be set.

* `custom_team_role_id` - (Optional) The ID of a custom role to assign to users. Conflicts with `standard_team_role`. One of `standard_team_role` or `custom_team_role_id` must be set.

* `is_admin` - (Optional) Whether group members should be Sysdig administrators. Default: `false`.

* `team_map` - (Required) Block defining team mapping. Maximum 1 block.

* `weight` - (Optional) Priority weight for conflict resolution. Lower numbers have higher priority. Must be between 1 and 32767. Default: `32767`.

### team_map

* `is_for_all_teams` - (Required) Whether the mapping applies to all teams.

* `team_ids` - (Optional) List of team IDs. Required when `is_for_all_teams` is `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the SSO group mapping.

## Import

SSO group mapping can be imported using the ID:

```
$ terraform import sysdig_sso_group_mapping.example 12345
```
