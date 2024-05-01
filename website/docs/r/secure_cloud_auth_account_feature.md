---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_cloud_auth_account_feature"
description: |-
  Creates a Sysdig Secure Cloud Account Feature using Cloudauth APIs.
---

# Resource: sysdig_secure_cloud_auth_account_feature

Creates a Sysdig Secure Cloud Account Feature using Cloudauth APIs.

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
  instance                   = "secure-scanning"
  service_principal_metadata = jsonencode({
	  gcp = {
      workload_identity_federation = {
          pool_provider_id = "some-pool-provider-id"
      }
      email = "some-service-account-email"
	  }
  })
}

resource "sysdig_secure_cloud_auth_account_feature" "sample" {
  account_id		             = sysdig_secure_cloud_auth_account.sample.id
  type                       = "FEATURE_SECURE_AGENTLESS_SCANNING"
  enabled                    = true
  components                 = ["COMPONENT_SERVICE_PRINCIPAL/secure-scanning"]
  flags                      = {
      "SCANNING_HOST_CONTAINER_ENABLED": "true"
  }
  depends_on = [ sysdig_secure_cloud_auth_account_component.sample ]
}
```

## Argument Reference

* `account_id` - (Required) Cloud Account created using resource sysdig_secure_cloud_auth_account.

* `type` - (Required) The type of feature to be created/added. e.g. `FEATURE_SECURE_CONFIG_POSTURE`.

* `enabled` - (Required) Whether or not to enable this feature on the given cloud account.

* `components` - (Required) Based on the feature type to be created, this is the list of components to be enabled on the cloud account.

* `flags` - (Optional) Based on the feature type to be created, these are the flags to be added to the feature on the cloud account.

-> **Note:** Please refer to Sysdig Secure API Documentation for the Cloud Accounts API for `feature` types and their related `components`.

-> **Note:** Since creation of component resource updates the account resource in the backend, in these configurations we indicate to Terraform to ignore `component` & `feature` attributes when planning updates to the remote account resource object.

## Attributes Reference

No additional attributes are exported.