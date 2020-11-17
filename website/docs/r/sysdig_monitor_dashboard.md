---
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_dashboard"
sidebar_current: "docs-sysdig-monitor-dashboard"
description: |-
  Creates a Sysdig Monitor Dashboard.
---

# sysdig\_monitor\_dashboard

Creates a Sysdig Monitor Dashboard using PromQL queries.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_monitor_dashboard" "dashboard" {
	name = "Example Dashboard"
	description = "Example Dashboard description"

	panel {
		pos_x = 0
		pos_y = 0
		width = 12 # Maximum size: 24
		height = 6
		type = "timechart" # timechart or number
		name = "Example panel"
		description = "Description"

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "percent"
		}
		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "number"
		}
	}

	panel {
		pos_x = 12
		pos_y = 0
		width = 12
		height = 6
		type = "number"
		name = "example panel - 2"
		description = "description of panel 2"

		query {
			promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
			unit = "time"
		}
	}
}
```

## Argument Reference


* `name` - (Required) The name of the Dashboard.

* `description` - (Optional) Description of the dashboard.

* `public` - (Optional) Define if the dashboard can be accessible without requiring the user to be logged in.

* `panel` - (Required) At least 1 panel is required to define a Dashboard.

### panel

The whole screen for a dashboard is separated in 24 squares of width. All the panels must not
overlap with other panels.
For example, if you position a panel in x: 0, y: 0, and you give it a width of 12, 
then you can position another panel in x: 12, y: 0 with a width of 12.

The following arguments are supported:

* `pos_x` - (Required) Position of the panel in the X axis. Min value: 0, max value: 23.

* `pos_y` - (Required) Position of the panel in the Y axis. Min value: 0.

* `width` - (Required) Width of the panel. Min value: 1, max value: 24. 

* `height` - (Required) Height of the panel. Min value: 1.

* `name` - (Required) Name of the panel.

* `description` - (Optional) Description of the panel.

* `type` - (Required) Kind of panel, must be either `timechart` or `number`.

* `query` - (Required) The PromQL query that will show information in the panel. 
            If the type of the panel is `timechart`, then it can be specified multiple 
            times, to have multiple metrics in the same graph.
            If the type of the panel is `number` then only one can be specified.


### query 

The following arguments are supported:

* `promql` - (Required) The PromQL query. Must be a valid PromQL query with existing
             metrics in Sysdig Monitor.
             
* `unit` - (Required) The type of metric for this query. Can be one of: `percent`, `data`, `data rate`, 
            `number`, `number rate`, `time`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `public_token` - (Computed) Token defined when the dashboard is set Public.

* `version` - (Computed)  The current version of the Dashboard.
