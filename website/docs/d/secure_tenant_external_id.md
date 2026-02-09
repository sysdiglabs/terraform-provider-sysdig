---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_tenant_external_id"
description: |-
  Retrieves information about the Sysdig Secure Tenant External ID
---

# Data Source: sysdig_secure_tenant_external_id

Retrieves the **cloud onboarding** external ID for the Sysdig Secure Tenant. This ID is used when configuring trusted relationships for cloud account onboarding (e.g., AWS IAM role trust policies).

~> **Note:** This is *not* the Customer External ID shown in [Customer ID, Name, and External ID](https://docs.sysdig.com/en/administration/find-your-customer-id-and-name/). For the customer-level external ID, use [`sysdig_current_user`](current_user.md) and its `customer_external_id` attribute instead.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_tenant_external_id" "external_id" {}
```

## Argument Reference

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `external_id` - The cloud onboarding external ID for the Sysdig Secure Tenant.

