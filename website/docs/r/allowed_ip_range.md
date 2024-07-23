---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_allowed_ip_range"
description: |-
  Creates allowed IP range in Sysdig which can be used to restrict access to the Sysdig platform.
---

# Resource: sysdig_allowed_ip_range

Creates allowed IP range in Sysdig which can be used to restrict access to the Sysdig platform.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_allowed_ip_range" "example" {
  ip_range = "192.168.100.0/24"
  note     = "Office IP range"
  enabled  = true
}

```
This example creates an allowed IP range for 192.168.100.0/24, with a note indicating it's for an office IP range, and it's enabled.


## Argument Reference

* `ip_range` - (Required) The IP range to allow access to the Sysdig platform. Must be in CIDR notation.
* `enabled` - (Required) Specifies whether the IP range is enabled.
* `note` - (Optional) A note describing the allowed IP range.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:  
* `id` - The ID of the allowed IP range.

## Import

Sysdig allowed IP ranges can be imported using the ID, e.g.

```
$ terraform import sysdig_allowed_ip_range.example 12345
```
