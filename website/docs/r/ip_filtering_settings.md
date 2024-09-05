---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_ip_filtering_settings"
description: |-
  Creates a IP filters settings in Sysdig.
---

# Resource: sysdig_ip_filtering_settings

Configures settings for IP filters (`sysdig_ip_filter` resource) which can be used to restrict access to the Sysdig platform.
Currently, there is only one setting available: `ip_filtering_enabled` which enables or disables the IP filtering feature. To enable the feature, at least one IP range must be defined in the `sysdig_ip_filter` resource.

> **Warning**
> This resource is global and is allowed to have only one instance per customer.
> Please verify that all IP ranges are created before enabling the feature. Failure to include your IP range will block your access to Sysdig until you connect from an approved IP range.


The `sysdig_ip_filtering_settings` behaves differently from normal resources, in that Terraform does not destroy this resource.
On resource destruction, Terraform performs no actions in Sysdig.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_ip_filtering_settings" "example" {
  ip_filtering_enabled = true
}

```
This example enables the IP filtering feature.

## Argument Reference

* `ip_filtering_enabled` - (Required) Specifies whether the IP filtering feature is enabled.

## Attributes Reference

No additional attributes are exported.

## Import

Sysdig IP filters settings can be imported, e.g.

```
$ terraform import sysdig_ip_filtering_settings.example ip_filtering_settings_id
```
