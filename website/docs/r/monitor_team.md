---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_team"
description: |-
  Creates a Sysdig Monitor Team.
---

# Resource: sysdig_monitor_team

Creates a Sysdig Monitor Team.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_team" "devops" {
  name = "Monitoring DevOps team"

  entrypoint {
	type = "Explore"
  }
  
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

* `name` - (Required) The name of the Monitor Team. It must be unique and must not exist in Secure.

* `entrypoint` - (Required) Main entry point for the current team in the product. 
                 See the Entrypoint argument reference section for more information.

* `description` - (Optional) A description of the team.

* `theme` - (Optional) Colour of the team. Default: "#73A1F7".

* `scope_by` - (Optional) Scope for the team. Default: "container".

* `filter` - (Optional) If the team can only see some resources, 
             write down a filter of such resources.
             
* `use_sysdig_capture` - (Optional) Defines if the team is able to create Sysdig Capture files. 
                         Default: true.
                         
* `can_see_infrastructure_events` - (Optional) TODO. Default: false.

* `can_use_aws_data` - (Optional) TODO. Default: false.

* `user_roles` - (Optional) Multiple user roles can be specified.
                 Administrators of the account will be automatically added
                 to every new created team, so they don't need to be added as a
                 resource in the Terraform manifest.

### Entrypoint Argument Reference

* `type` - (Required) Main entrypoint for the team.
                      Valid options are: Explore, Dashboards, Events, Alerts, Settings.

* `selection` - (Optional) Sets up the defined Dashboard name as entrypoint.
                Warning: This field must only be added if the `type` is "Dashboards".

### User Role Argument Reference

* `email` - (Required) The email of the user in the group.

* `role` - (Optional) The role for the user in this group.
           Valid roles are: ROLE_TEAM_STANDARD, ROLE_TEAM_EDIT, ROLE_TEAM_READ, ROLE_TEAM_MANAGER or CustomRole ID.
           Default: ROLE_TEAM_STANDARD.
#### Custom Role example


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_team` - (Optional) Mark team as default team. Users with no designated team will be added to this team by default.

### IBM Cloud Monitoring arguments

* `enable_ibm_platform_metrics` - (Optional) Enable platform metrics on IBM Cloud Monitoring.

* `ibm_platform_metrics` - (Optional) Define platform metrics on IBM Cloud Monitoring.

## Import

Monitor Teams can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_team.example 12345
```
