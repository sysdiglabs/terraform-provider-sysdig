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

```terraform
resource "sysdig_group_mapping" "my_group" {
  group_name = "my-group"
  role = "ROLE_TEAM_STANDARD"

  team_map {
    all_teams = false
    team_ids = [sysdig_secure_team.my_team.id, sysdig_monitor_team.my_team.id]
  }
}

```

## Argument Reference

* `group_name` - (Required) The group name to be mapped.

* `role` - (Required) The role that is assigned to the users. It can be a standard role or a custom team role ID.

* `team_map` - (Required) Block to define team mapping.

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
