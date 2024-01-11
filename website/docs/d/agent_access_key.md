---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_agent_access_key"
description: |-
  Retrieves information about a agent access key from the access key id.
---

# Data Source: sysdig_custom_role

Retrieves information about a agent access key from the access key id.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_agent_access_key" "sysdig_agent_access_key" {
    agent_key = "abcke91d-2495-4192-8721-0a6bae7deec9"
}
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `metadata` - The metadata of agent access key.

* `team_id` - The team id of the agent access key.

* `team_name` - The team name of the agent access key.

* `limit` - The limit of the agent access key.

* `reservation` - The reservation of the agent access key.

* `agents_connected` - The number of agents connected with that agent access key.

* `enabled` - Whether the agent access key is enabled or not.


