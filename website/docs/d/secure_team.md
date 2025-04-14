---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_team"
description: |-
  Retrieves information about a specific secure team in Sysdig
---

# sysdig_secure_team

The `sysdig_secure_team` data source retrieves information about a specific secure team in Sysdig.

## Example Usage

```terraform
data "sysdig_secure_team" "example" {
  id = "812371"
}
```

## Argument Reference

- `id` - (Required) The ID of the secure team.

## Attribute Reference

- `name` - The name of the secure team.
- `description` - The description of the secure team.
- `filter` - The filter applied to the team.
- `scope_by` - The scope of the team.
- `use_sysdig_capture` - Whether the team can use Sysdig capture.
- `default_team` - Whether the team is the default team.
- `user_roles` - The roles assigned to users in the team.
- `zone_ids` - The IDs of the zones associated with the team.
- `all_zones` - Whether the team has access to all zones.
