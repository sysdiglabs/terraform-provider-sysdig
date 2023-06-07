---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_team"
description: |-
Creates a Sysdig Monitor Team.
---

# Resource: sysdig_monitor_team

Creates a Sysdig Monitor Team.

## Example Usage

```terraform
resource "sysdig_monitor_team" "devops" {
  name = "Monitoring DevOps team"
}
```

## Argument Reference

* `name` - (Required) The name of the Monitor Team. It must be unique and must not exist in Secure.
