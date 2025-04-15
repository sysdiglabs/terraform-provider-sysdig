---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_teams"
description: |-
  Retrieves information about a specific monitor teams in Sysdig
---

# sysdig_monitor_teams

The `sysdig_monitor_teams` data source retrieves a list of all monitor teams in Sysdig.

## Example Usage

```terraform
data "sysdig_monitor_teams" "example" {}
```

## Attribute Reference

- `teams` - A list of monitor teams. Each team has the following attributes:
  - `id` - The ID of the monitor team.
  - `name` - The name of the monitor team.
  - `description` - The description of the monitor team.
  - `filter` - The filter applied to the team.
  - `scope_by` - The scope of the team.
  - `can_use_sysdig_capture` - Whether the team can use Sysdig capture.
  - `can_see_infrastructure_events` - Whether the team can see infrastructure events.
  - `can_use_aws_data` - Whether the team can use AWS data.
  - `default_team` - Whether the team is the default team.
  - `user_roles` - The roles assigned to users in the team.
  - `version` - The version of the monitor team.
  - `theme` - The theme of the monitor team.
