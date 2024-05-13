---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_cloud_auth_account_component"
description: |-
  Creates a Sysdig Secure Cloud Account Component using Cloudauth APIs.
---

# Resource: sysdig_secure_cloud_auth_account_component

Creates a Sysdig Secure Cloud Account Component using Cloudauth APIs.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id        = "mygcpproject"
  provider_type      = "PROVIDER_GCP"
  enabled            = true
  lifecycle {
	  ignore_changes = [
	    component,
	    feature
	  ]
  }
}
resource "sysdig_secure_cloud_auth_account_component" "sample" {
  account_id		             = sysdig_secure_cloud_auth_account.sample.id
  type                       = "COMPONENT_SERVICE_PRINCIPAL"
  instance                   = "secure-posture"
  service_principal_metadata = jsonencode({
	  gcp = {
		  key = "gcp-sa-key"
	  }
  })
}
```

## Argument Reference

* `account_id` - (Required) Cloud Account created using resource sysdig_secure_cloud_auth_account.

* `type` - (Required) The type of component to be created. e.g. `COMPONENT_SERVICE_PRINCIPAL`.

* `instance` - (Required) The component instance to be created, identified by a specific string. e.g. `secure-posture`, `secure-runtime`, etc.

* `<component>_metadata` - (Optional) Based on the component type created, this is the metadata information passed to enable the component on the account.

-> **Note:** Please refer to Sysdig Secure API Documentation for the Cloud Accounts API for metadata types for `component`.

-> **Note:** Since creation of component resource updates the account resource in the backend, in these configurations we indicate to Terraform to ignore `component` & `feature` attributes when planning updates to the remote account resource object.

## Attributes Reference

No additional attributes are exported.