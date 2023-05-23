---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_falco"
description: |-
  Retrieves a Sysdig Secure Falco Rule.
---

# Data Source: sysdig_secure_rule_falco

Retrieves the information of an existing Sysdig Secure Falco Rule.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_rule_falco" "example" {
    name = "Terminal shell in container"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure rule to retrieve.
* `source` - (Optional) The source of the Secure rule to retrieve. This is required if a rule with this name exists in
  multiple sources.
* `index` - (Optional) The index of the Secure rule to retrieve in the event of multiple rules. The default value is 0.
  see section on rules with appended rules.

## Attributes Reference

In addition to the argument above, the following attributes are exported:

* `description` - The description of Secure rule.
* `tags` - A list of tags for this rule.
* `condition` - A [Falco condition](https://falco.org/docs/rules/) is simply a Boolean predicate on Sysdig events
  expressed using the Sysdig [filter syntax](http://www.sysdig.org/wiki/sysdig-user-guide/#filtering) and macro terms.
* `output` - Add additional information to each Falco notification's output.
* `priority` - The priority of the Falco rule. It can be: "emergency", "alert", "critical", "error", "warning", "notice", "info" or "debug". By default is "warning".
* `exceptions` - The exceptions key is a list of identifier plus list of tuples of filtercheck fields. See below for details.
* `append` - This indicates that the rule being created appends the condition to an existing Sysdig-provided rule
* `version` - Current version of the resource in Sysdig Secure.

### Exceptions

Starting in 0.28.0, Falco supports an optional exceptions property to rules. The exceptions key is a list of identifier plus list of tuples of filtercheck fields.
For more information about the syntax of the exceptions, check the [official Falco documentation](https://falco.org/docs/rules/exceptions/).

Supported fields for exceptions:

* `name` - The name of the exception.
* `fields` - Contains one or more fields that will extract a value from the syscall/k8s_audit events.
* `comps` - Contains comparison operators that align 1-1 with the items in the fields property.
* `values` - Contains tuples of values. Each item in the tuple should align 1-1 with the corresponding field
  and comparison operator. 

## Rules with Appended Rules

In the event that a rule has appended rules, the data source can return the default rule and its appended rules using
the `index` argument. This can be combined with the `sysdig_secure_rule_falco_count` data source to easily retrieve all
of the rules in the rule group.

An example of how this could be used follows:

```terraform
data "sysdig_secure_rule_falco_count" "disallowed_container" {
  name = "Launch Disallowed Container"
  source = "syscall"
}

data "sysdig_secure_rule_falco" "disallowed_container" {
  count = data.sysdig_secure_rule_falco_count.disallowed_container.rule_count
  name = "Launch Disallowed Container"
  source = "syscall"
  index = "${count.index}"
  depends_on = [ data.sysdig_secure_rule_falco_count.disallowed_container ]
}

output "disallowed_container_rule_group" {
  value = ["${data.sysdig_secure_rule_falco.disallowed_container.*}"]
}
```
