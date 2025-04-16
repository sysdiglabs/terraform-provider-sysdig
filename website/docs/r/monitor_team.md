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
	  type = "DashboardTemplates"
    selection = "view.net.http"
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

  filter = "kubernetes.namespace.name in (\"kube-system\") and kubernetes.deployment.name in (\"coredns\")"
  prometheus_remote_write_metrics_filter = "kube_cluster_name in (\"test-cluster\", \"test-k8s-data\") and kube_deployment_name  = \"coredns\" and my_metric starts with \"prefix\" and not my_metric contains \"prefix-test\""
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

* `scope_by` - (Optional) Scope for the team, either "container" or "host". Default: "container". If set to host, Team members can see all Host-level and Container-level information. If set to Container, Team members can see only Container-level information.

* `filter` - (Optional) Use this option to select which Agent Metrics data users of this team can view. Not setting it will allow users to see all Agent Metrics data.

* `prometheus_remote_write_metrics_filter` - (Optional) Use this option to select which Prometheus Remote Write data users of this team can view. Not setting it will allow users to see all Prometheus Remote Write data.

* `use_sysdig_capture` - (Optional) Defines if the team is able to create Sysdig Capture files.  Default: true.

* `can_see_infrastructure_events` - (Optional) Enable this option to allow this team to view all Infrastructure and Custom Events from every user and agent. Otherwise, this team will only see infrastructure events sent specifically to this team. Default: false.

* `can_use_aws_data` - (Optional) Enable this option to give this team access to AWS metrics and tags. All AWS data is made available, regardless of the team’s Scope. Default: false.

* `can_use_agent_cli` - (Optional) Enable this option to give this team access to Using the Agent Console. Default: true.

* `user_roles` - (Optional) Multiple user roles can be specified.
                 Administrators of the account will be automatically added
                 to every new created team, so they don't need to be added as a
                 resource in the Terraform manifest.

### Entrypoint Argument Reference

* `type` - (Required) Main entrypoint for the team.
                      Valid options are: `Explore`, `Dashboards`, `Events`, `Alerts`, `Settings`, `DashboardTemplates`, `Overview`.

* `selection` - (Optional) Sets up the defined Dashboard name as entrypoint.
                Warning: This field must only be added if the `type` is `Dashboards`, and the value is the numeric id of the selected dashboard, or `DashboardTemplates`, and the value is the id (dotted name) of the selected dashboard template.

### User Role Argument Reference

* `email` - (Required) The email of the user in the group.

* `role` - (Optional) The role for the user in this group.
           Valid roles are: ROLE_TEAM_STANDARD, ROLE_TEAM_EDIT, ROLE_TEAM_READ, ROLE_TEAM_MANAGER or CustomRole ID.<br/>
           Default: ROLE_TEAM_STANDARD.<br/>
           Note: CustomRole ID can be referenced from `sysdig_custom_role` resource or `sysdig_custom_role` data source

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default_team` - (Optional) Mark team as default team. Users with no designated team will be added to this team by default.

### IBM Cloud Monitoring arguments

* `enable_ibm_platform_metrics` - (Optional) Enable Platform Metrics on IBM Cloud Monitoring.

* `ibm_platform_metrics` - (Optional) Use this option to select which Platform Metrics data users of this team can view. Not setting it will allow users to see all Platform Metrics data.

## Import

Monitor Teams can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_team.example 12345
```
