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
