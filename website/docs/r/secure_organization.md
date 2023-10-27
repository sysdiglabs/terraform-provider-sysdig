---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_organization"
description: |- 
  Creates a Sysdig Secure Organization 
---

# Resource: sysdig_secure_organization

Creates a Sysdig Secure Organization.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "google_cloud_project_id"
  provider_type = "PROVIDER_GCP"
  enabled       = "true"
}
resource "sysdig_secure_organization" "sample" {
  management_account_id	    = sysdig_secure_cloud_auth_account.sample.id 
}
```

## Argument Reference

* `management_account_id` - (Required) Cloud Account created using resource sysdig_secure_cloud_auth_account.
* `organizational_unit_ids` - (Optional) List of organizational unit identifiers from which to onboard. If empty, the entire organization is onboarded. 

## Attributes Reference

No additional attributes are exported.
