---
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_dashboard"
description: |-
  Creates a Sysdig Monitor Dashboard.
---

# Resource: sysdig\_monitor\_dashboard

Creates a Sysdig Monitor Dashboard using PromQL queries.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_monitor_dashboard" "dashboard" {
  name        = "Example Dashboard"
  description = "Example Dashboard description"

  scope {
    metric     = "kubernetes.cluster.name"
    comparator = "in"
    value      = ["prod", "dev"]
    variable   = "cluster_name"
  }

  scope {
    metric   = "host.hostName"
    variable = "hostname"
  }

  panel {
    pos_x       = 0
    pos_y       = 0
    width       = 12 # Maximum size: 24
    height      = 6
    type        = "timechart" # timechart or number
    name        = "Example panel"
    description = "Description"

    query {
      promql = "avg_over_time(sysdig_host_cpu_used_percent{host_name=$hostname}[$__interval])"
      unit   = "percent"
    }
    query {
      promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
      unit   = "number"
    }
  }

  panel {
    pos_x       = 12
    pos_y       = 0
    width       = 12
    height      = 6
    type        = "number"
    name        = "example panel - 2"
    description = "description of panel 2"

    query {
      promql = "avg(avg_over_time(sysdig_host_cpu_used_percent[$__interval]))"
      unit   = "time"
    }
  }

  panel {
    pos_x                  = 12
    pos_y                  = 12
    width                  = 12
    height                 = 6
    type                   = "text"
    name                   = "example panel - 2"
    content                = "description of panel 2"
    visible_title          = true
    autosize_text          = true
    transparent_background = true
  }
}

```

## Argument Reference


* `name` - (Required) The name of the Dashboard.

* `description` - (Optional) Description of the dashboard.

* `public` - (Optional) Define if the dashboard can be accessible without requiring the user to be logged in.
  
* `scope` - (Optional) Define the scope of the dashboard and variables for these metrics.

* `panel` - (Required) At least 1 panel is required to define a Dashboard.


### scope

Dashboard scope defines what data is valid for aggregation and display within the dashboard.
See more info about how to [use the scope in a PromQL query](https://docs.sysdig.com/en/using-promql.html#UUID-2314cf2d-3466-d7a5-142a-30a9e63053d0_UUID-8dfed5eb-8c48-8f94-4e3a-61b051fb9b440) in the official documentation.

The following arguments are supported to configure a scope:

* `metric` - (Required) Metric to scope by, common examples are `host.hostName`, `kubernetes.namespace.name` or `kubernetes.cluster.name`, but you can use all the Sysdig-supported values shown in the UI. Note that kubernetes-related values only appear when Sysdig detects Kubernetes metadata.

* `comparator` - (Optional) Operator to relate the metric with some value. It is only required if the value to filter by is set, or the variable field is not set. Valid values are: `in`, `notIn`, `equals`, `notEquals`, `contains`, `notContains` and `startsWith`.

* `value` - (Optional) List of values to filter by, if comparator is set. If the comparator is not `in` or `notIn` the list must contain only 1 value.
  
* `variable` - (Optional) Assigns this metric to a value name and allows PromQL to reference it.


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

* `type` - (Required) Kind of panel, must be either `timechart`, `number` or `text`.

* `query` - (Optional) The PromQL query that will show information in the panel. 
            If the type of the panel is `timechart`, then it can be specified multiple 
            times, to have multiple metrics in the same graph.
            If the type of the panel is `number` then only one can be specified.
            This field is required if the panel type is `timechart` or `number`.

* `content` - (Optional) This field is required if the panel type is `text`. It represents the 
               text that will be displayed in the panel.

* `visible_title` - (Optional) If true, the title of the panel will be displayed. Default: false.
                    This field is ignored for all panel types except `text`.

* `autosize_text` - (Optional) If true, the text will be autosized in the panel.
                    This field is ignored for all panel types except `text`.
  
* `transparent_background` - (Optional) If true, the panel will have a transparent background.
                             This field is ignored for all panel types except `text`.

### query 

To scope a panel built from a PromQL query, you must use a scope variable within the query. The variable will take the value of the referenced scope parameter, and the PromQL panel will change accordingly.
There are two predefined variables available:

- `$__interval` represents the time interval defined based on the time range. This will help to adapt the time range for different operations, such as rate and avg_over_time, and prevent displaying empty graphs due to the change in the granularity of the data.

- `$__range` represents the time interval defined for the dashboard. This is used to adapt operations like calculating average for a time frame selected.

The following arguments are supported:

* `promql` - (Required) The PromQL query. Must be a valid PromQL query with existing
             metrics in Sysdig Monitor.
             
* `unit` - (Required) The type of metric for this query. Can be one of: `percent`, `data`, `data rate`, 
            `number`, `number rate`, `time`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `public_token` - (Computed) Token defined when the dashboard is set Public.

* `version` - (Computed)  The current version of the Dashboard.

## Import

Monitor dashboards can be imported using the dashboard ID, e.g.

```
$ terraform import sysdig_monitor_dashboard.example 12345
```

Only dashboards that contain supported panels can be imported. Currently supported panel types are:
- PromQL timecharts
- PromQL numbers
- Text

Only dashboards that contain supported query types can be imported. Currently supported query types:
- Percent
- Data
- Data rate
- Number
- Number rate
- Time