---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_tenant_external_id"
description: |-
  Retrieves information about the Sysdig Secure Tenant External ID
---

# Data Source: sysdig_secure_tenant_external_id

Retrieves information about the Sysdig Secure Tenant External ID

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_tenant_external_id" "external_id" {}
```

## Argument Reference

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `external_id` - String identifier for external id value

