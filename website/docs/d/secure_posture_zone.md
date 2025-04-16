---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_posture_zone"
description: |-
  Retrieves Posture Zone by ID.
---

# sysdig_secure_posture_zone Data Source

The `sysdig_secure_posture_zone` data source allows you to retrieve information about a specific secure posture zone by its ID.

## Example Usage

```terraform
data "sysdig_secure_posture_zone" "example" {
  id = "454678"
}
```

## Argument Reference

The following arguments are supported:

- `id` (Required) - The ID of the secure posture zone to retrieve.

## Attribute Reference

The following attributes are exported:

- `name` - The name of the secure posture zone.
- `description` - The description of the secure posture zone.
- `policy_ids` - A list of policy IDs associated with the secure posture zone.
- `author` - The author of the secure posture zone.
- `last_modified_by` - The user who last modified the secure posture zone.
- `last_updated` - The timestamp of the last update to the secure posture zone.
- `scopes` - A list of scopes associated with the secure posture zone. Each scope contains:
  - `target_type` - The target type of the scope.
  - `rules` - The rules associated with the scope.
