---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_network"
description: |-
  Retrieves a Sysdig Secure Network Rule.
---

# Data Source: sysdig_secure_rule_network

Retrieves the information of an existing Sysdig Secure Network Rule.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_rule_network" "example" {
    name = "Disallowed SSH Connection"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure rule to retrieve.

## Attributes Reference

In addition to the argument above, the following attributes are exported:

* `description` - The description of Secure rule.
* `tags` - A list of tags for this rule.
* `block_inbound` - Detect if there is an inbound connection.
* `block_outbound` - Detect if there is an outbound connection.
* `tcp` - A block with the properties `matching` and `ports` for TCP connections.
* `udp` - A block with the properties `matching` and `ports` for UDP connections.
* `version` - Current version of the resource in Sysdig Secure.

## Connection Blocks

The `tcp` and `udp` blocks will have the the following attributes:

* `matching` - Defines if the port matches or not with the provided list.
* `ports` - List of ports to match.
