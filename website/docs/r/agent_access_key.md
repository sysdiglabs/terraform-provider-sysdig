---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_agent_access_key"
description: |-
  Retrieves information about a agent access key from the access key id.
---

# Resource: sysdig_agent_access_key

Retrieves information about an agent access key from the access key id.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_agent_access_key" "my_agent_access_key" {
  limit = 11
  reservation = 1
  team_id = 50012099
  metadata = {
    "test" = "yes"
    "environment" = "development"
    "team" = "awesome_team"
  }
  enabled = true
}
```

## Argument Reference

In addition to all arguments above, the following attributes are exported:

* `metadata` - (Optional) The metadata of agent access key.

* `team_id` - (Optional) The team id of the agent access key.

* `limit` - (Optional) The limit of the agent access key.

* `reservation` - (Optional) The reservation of the agent access key.

* `enabled` - (Optional) Whether the agent access key is enabled or not. It is only used in update actions, an agent access keys can be deleted only if it's disabled.

## Attributes Reference

* `access_key` - The agent access key.

* `date_disabled` - Date when the agent key was last disabled.

* `date_created` - Date when the agent key was created.

## Import

Sysdig group mapping can be imported using the ID, e.g.

```
$ terraform import sysdig_agent_access_key.my_agent_access_key "631123"
```
