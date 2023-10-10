---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_team_service_account"
description: |-
  Creates a team service account in Sysdig.
---

# Resource: sysdig_team_service_account

Creates a team service account in Sysdig.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "time_static" "example" {
  rfc3339 = "2025-01-01T00:00:00Z"
}

resource "sysdig_monitor_team" "devops" {
  name = "Monitoring DevOps team"

  entrypoint {
    type = "Explore"
  }
}

resource "sysdig_team_service_account" "service-account" {
  name = "read only"
  role = "ROLE_TEAM_READ"
  expiration_date = time_static.example.unix
  team_id = sysdig_monitor_team.devops.id
}

```

## Argument Reference

* `name` - (Required) The team service account name.

* `role` - (Required) The role that is assigned to the service account. It can be a standard role or a custom team role ID.

* `expiration_date` (Required) The service account expiration date.

* `team_id` - (Required) The team where the service account belongs to.

* `system_role`- The service account system role. The only value supported is `ROLE_SERVICE_ACCOUNT`


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `date_created` - The team service account creation date

* `api_key` - The api key to be using in API calls

## Import

Sysdig team service account can be imported using the ID, e.g.

```
$ terraform import sysdig_team_service_account.my_team_service_account 10
```
