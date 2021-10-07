---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_network"
description: |-
  Creates a Sysdig Secure Network Rule.
---

# Resource: sysdig_secure_rule_network

Creates a Sysdig Secure Network Rule.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_secure_rule_network" "example" {
  name = "Disallowed SSH Connection"
  description = "Detect any new ssh connection to a host other than those in an allowed group of hosts"
  tags = ["network", "mitre_remote_service"]

  block_inbound = true
  block_outbound = true

  tcp {
    matching = true // default
    ports = [22]
  }

  udp {
    matching = true // default
    ports = [22]
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Optional) The description of Secure rule. By default is empty.
* `tags` - (Optional) A list of tags for this rule.

### Disallow incoming or outgoing connections

* `block_inbound` - (Required) Detect if there is an inbound connection.
* `block_outbound` - (Required) Detect if there is an outbound connection.

### Detect TCP Connections

* `matching` - (Optional) Defines if the port matches or not with the provided list. Default is true.
* `ports` - (Required) List of ports to match.

### Detect UDP Connections

* `matching` - (Optional) Defines if the port matches or not with the provided list. Default is true.
* `ports` - (Required) List of ports to match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

## Import

Secure network runtime rules can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_rule_network.example 12345
```