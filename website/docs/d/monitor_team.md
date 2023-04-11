---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_team"
description: |-
Retrieves information about a Monitor team
---

# Data Source: sysdig_monitor_team

Retrieves information about a Monitor team

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_monitor_team" "monitor_by_name"{
  name = "team name"
}
```


## Argument Reference

* `name` - (Required) The name of the Team to be retrieved.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `entrypoint` - Main entry point for the current team in the product.
* `description` - A description of the team.
* `theme` - Colour of the team.
* `scope_by` - Scope for the team.
* `filter` - Filters on the team.
* `use_sysdig_capture` -  Defines if the team is able to create Sysdig Capture files.
* `can_see_infrastructure_events` - Defines if the team is able to use infrastructure events.
* `can_use_aws_data` - Defines if the team is able to use AWS data.
* `user_roles` - User roles in the team.
* `default_team` - Flag which indicates if team is default.

### IBM Cloud Monitoring attributes
* `enable_ibm_platform_metrics` - Flag which indicates if platform metrics are enabled.
* `ibm_platform_metrics` - Defined platform metrics on the team.

