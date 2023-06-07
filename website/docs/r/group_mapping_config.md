---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_group_mapping_config"
description: |-
  Sets the group mapping conflicts resolution strategies in Sysdig.
---

# Resource: sysdig_group_mapping_config

Sets the group mapping conflicts resolution strategies in Sysdig.

> **Warning**
> This resource is global and is allowed to have only one configuration per customer

The `sysdig_group_mapping_config` behaves differently from normal resources, in that Terraform does not destroy this resource.
On resource destruction, Terraform performs no actions in Sysdig.

## Example Usage

```terraform
resource "sysdig_group_mapping_config" "resolution_strategies" {
  no_mapping_strategy = "UNAUTHORIZED"
  different_team_same_role_strategy = "UNAUTHORIZED"
}
```

## Argument Reference

* `no_mapping_strategy` - (Required) Sets how the system behaves when no group mapping information received from the IdP or Group information received, but the user is not a member of any mapped group. Possible values are: `UNAUTHORIZED`, `DEFAULT_TEAM_DEFAULT_ROLE`

* `different_team_same_role_strategy` - (Required) Sets how the system behaves when conflicting group mapping information received. Possible values are: `UNAUTHORIZED`, `FIRST_MATCH`, `WEIGHTED`

## Attributes Reference

No additional attributes are exported.

## Import

Sysdig group mapping config can be imported, e.g.

```
$ terraform import sysdig_group_mapping_config.resolution_strategies conflicts_resolution_strategies
```
