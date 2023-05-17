---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_falco_count"
description: |-
  Retrieves the count of rules (including appends) for a named falco rule.
---

# Data Source: sysdig_secure_rule_falco_count

Retrieves the count of rules (including appends) for a named falco rule.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_rule_falco_count" "example" {
    name = "Terminal shell in container"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure rule to retrieve.
* `source` - (Optional) The source of the Secure rule to retrieve. This is required if a rule with this name exists in
  multiple sources.

## Attributes Reference

In addition to the argument above, the following attributes are exported:

* `rule_count` - The number of rules (including appends).
