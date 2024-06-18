---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_inhibition_rule"
description: |-
  Creates a Sysdig Monitor Inhibition Rule.
---

# Resource: sysdig_monitor_inhibition_rule

Creates a Sysdig Monitor Inhibition Rule.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_inhibition_rule" "sample" {
  name = "Example Inhibition Rule"
  description = "Example description"
  enabled = true

  source_matchers {
    label_name = "alertname"
    operator = "EQUALS"
    value = "networkAlert"
  }

  source_matchers {
    label_name = "device_type"
    operator = "EQUALS"
    value = "firewall"
  }

  target_matchers {
    label_name = "device_type"
    operator = "REGEXP_MATCHES"
    value = ".*server.*"
  }

  equal = ["kube_cluster_name"]
}
```

## Argument Reference

* `name` - (Optional) The name of the Inhibition Rule. If provided, it must be unique.

* `description` - (Optional) The description of the Inhibition Rule.

* `enabled` - (Optional) Whether to enable the Inhibition Rule. Default: `true`.

* `equal` - (Optional) List of labels that must have an equal value in the source and target alert for the inhibition to take effect.

### `source_matchers`

List of source matchers for which one or more alerts have to exist for the inhibition to take effect.

It is a list of objects with the following fields:

* `label_name`: (Required) Label to match.

* `operator`: (Required) Match operator. It can be `EQUALS`, `NOT_EQUALS`, `REGEXP_MATCHES`, `NOT_REGEXP_MATCHES`.

* `value`: (Required) Label value to match in case operator is of type equality, or a valid regular expression in case of operator is of type regex.

### `target_matchers`

List of target matchers that have to be fulfilled by the target alerts to be muted.

It is a list of objects with the following fields:

* `label_name`: (Required) Label to match.

* `operator`: (Required) Match operator. It can be `EQUALS`, `NOT_EQUALS`, `REGEXP_MATCHES`, `NOT_REGEXP_MATCHES`.

* `value`: (Required) Label value to match in case `operator` is of type equality, or regular expression in case of `operator` is of type regex.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the Inhibition Rule.

* `version` - (Computed) The current version of the Inhibition Rule.

## Import

Inhibition Rules for Monitor can be imported using the ID, e.g.

```
$ terraform import sysdig_monitor_inhibition_rule.example 12345
```
