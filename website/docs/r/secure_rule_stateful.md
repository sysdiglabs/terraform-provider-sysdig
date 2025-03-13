---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_stateful"
description: |-
  Creates a Sysdig Secure Stateful Rule Append.
---

# Resource: sysdig_secure_rule_stateful

Creates a Sysdig Secure Stateful Rule Append.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_rule_stateful" "stateful_rule" {
  name = "API Gateway Enumeration Detected"
  source = "awscloudtrail_stateful"
  ruletype = "STATEFUL_SEQUENCE"
  exceptions {
      values = jsonencode([["user_abc", ["12345"]]])
      name = "user_accountid"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Stateful rule that the exception is being appended to.
* `source` - (Required) The source of the event. We currently support the "awscloudtrail_stateful" source.
* `exceptions` - (Required) The exceptions key is a list of identifier plus list of tuples of filtercheck fields. See below for details.
* `append` - (Optional) This indicates that the rule being created appends the condition to an existing Sysdig-provided. For stateful rules, the default value is true.
* `ruletype` - (Required) The type of Stateful rule being appended to. We currently support "STATEFUL_SEQUENCE", "STATEFUL_COUNT", and "STATEFUL_UNIQ_PERCENT".

### Exceptions
Supported fields for exceptions:

* `name` - (Required) The name of the exception.
* `values` - (Required) Contains tuples of values. Each item in the tuple should align 1-1 with the corresponding field
  and comparison operator. Since the value can be a string, a list of strings or a list of a list of strings, the value
  of this field must be supplied in JSON format. You can use the default `jsonencode` function to provide this value.
  See the usage example on the top.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

