---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_team"
description: |-
  Creates a Sysdig Secure Team.
---

# Resource: sysdig_secure_team

Creates a Sysdig Secure Team.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_team" "devops" {
  name = "DevOps team"
  
  user_roles {
    email = data.sysdig_current_user.me.email
    role = "ROLE_TEAM_MANAGER"
  }

  user_roles {
    email = "john.doe@example.com"
    role = "ROLE_TEAM_STANDARD"
  }

  user_roles {
    email = "john.smith@example.com"
    role = data.sysdig_custom_role.custom_role.id
  }
}
 
data "sysdig_current_user" "me" {
}

data "sysdig_custom_role" "custom_role" {
  name = "CustomRoleName"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure Team. It must be unique and must not exist in Monitor.

* `description` - (Optional) A description of the team.

* `theme` - (Optional) Colour of the team. Default: "#73A1F7".

* `scope_by` - (Optional) Scope for the team. Default: "container".

* `filter` - (Optional) If the team can only see some resources, 
             write down a filter of such resources.
             
* `use_sysdig_capture` - (Optional) Defines if the team is able to create Sysdig Capture files. 
                         Default: true.
                         
* `default_team` - (Optional) Defines if the team is the default one. Warning: only one can be the default,
                   if you define multiple default teams, Terraform will be updating the API in every execution,
                   even if the state hasn't changed.

* `user_roles` - (Optional) Multiple user roles can be specified.
                 Administrators of the account will be automatically added
                 to every new created team, so they don't need to be added as a
                 resource in the Terraform manifest.

* `zone_ids` - (Optional) List of zone IDs attached to the team. If `all_zones` is specified this argument needs to be omitted.

* `all_zones` - (Optional) Attach all zones to the team. If this argument is enabled then `zone_ids` needs to be omitted.
                         
### User Role Argument Reference

* `email` - (Required) The email of the user in the group.

* `role` - (Optional) The role for the user in this group.
           Valid roles are: ROLE_TEAM_STANDARD, ROLE_TEAM_EDIT, ROLE_TEAM_READ, ROLE_TEAM_MANAGER or CustomRole ID.
           Default: ROLE_TEAM_STANDARD.

## Attributes Reference

No additional attributes are exported.

### IBM Workload protection arguments

* `enable_ibm_platform_metrics` - (Optional) Enable platform metrics on IBM Cloud Monitoring.

* `ibm_platform_metrics` - (Optional) Define platform metrics on IBM Cloud Monitoring.

## Import

Secure Teams can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_team.example 12345
```
