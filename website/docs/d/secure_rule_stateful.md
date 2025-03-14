---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_stateful"
description: |-
  Retrieves a Sysdig Secure Stateful Rule.
---

# Data Source: sysdig_secure_rule_stateful

Retrieves the information of an existing Sysdig Secure Stateful Rule.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_rule_stateful" "example" {
    name = "Access Key Enumeration Detected"
    source = "awscloudtrail_stateful"
    ruletype = "STATEFUL_SEQUENCE"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure rule to retrieve.
* `source` - (Required) The source of the Secure rule to retrieve.
* `ruletype` - (Required) The type of the Secure rule to retrieve.

## Attributes Reference

In addition to the argument above, the following attributes are exported:

* `exceptions` - The exceptions key is a list of identifier plus list of tuples of filtercheck fields. See below for details.
* `append` - This indicates that the rule being created appends the condition to an existing Sysdig-provided rule

### Exceptions

Stateful rules support an optional exceptions property to rules. The exceptions key is a list of identifier plus list of tuples of filtercheck fields.

Supported fields for exceptions:

* `name` - The name of the existing exception definition.
* `values` - Contains tuples of values. Each item in the tuple should align 1-1 with the corresponding field
  and comparison operator. 
