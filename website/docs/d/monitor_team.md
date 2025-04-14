---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_team"
description: |-
  Retrieves information about a specific monitor team in Sysdig
---

# sysdig_monitor_team

The `sysdig_monitor_team` data source retrieves information about a specific monitor team in Sysdig.

## Example Usage

```terraform
data "sysdig_monitor_team" "example" {
  id = "812371"
}
```

## Argument Reference

- `id` - (Required) The ID of the monitor team.

## Attribute Reference

- `name` - The name of the monitor team.
- `description` - The description of the monitor team.
- `entrypoint` - The entrypoint configuration for the team.
- `filter` - The filter applied to the team.
- `scope_by` - The scope of the team.
- `can_use_sysdig_capture` - Whether the team can use Sysdig capture.
- `can_see_infrastructure_events` - Whether the team can see infrastructure events.
- `can_use_aws_data` - Whether the team can use AWS data.
- `default_team` - Whether the team is the default team.
- `user_roles` - The roles assigned to users in the team.
