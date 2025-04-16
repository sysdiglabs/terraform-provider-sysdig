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
