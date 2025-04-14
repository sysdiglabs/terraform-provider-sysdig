---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_teams"
description: |-
  Retrieves information about a specific secure teams in Sysdig
---

# sysdig_secure_teams

The `sysdig_secure_teams` data source retrieves a list of all secure teams in Sysdig.

## Example Usage

```terraform
data "sysdig_secure_teams" "example" {}
```

## Attribute Reference

- `teams` - A list of secure teams. Each team has the following attributes:
  - `id` - The ID of the secure team.
  - `name` - The name of the secure team.
  - `description` - The description of the secure team.
  - `filter` - The filter applied to the team.
  - `scope_by` - The scope of the team.
  - `use_sysdig_capture` - Whether the team can use Sysdig capture.
  - `default_team` - Whether the team is the default team.
  - `user_roles` - The roles assigned to users in the team.
  - `zone_ids` - The IDs of the zones associated with the team.
  - `all_zones` - Whether the team has access to all zones.
