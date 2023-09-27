---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_cloud_auth_account"
description: |-
  Creates a Sysdig Secure Cloud Account using Cloudauth APIs.
---

# Resource: sysdig_secure_cloud_auth_account

Creates a Sysdig Secure Cloud Account using Cloudauth APIs.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "mygcpproject"
  provider_type = "PROVIDER_GCP"
  enabled       = true
}
```

## Argument Reference

* `provider_id` - (Required) The unique identifier of the cloud account. e.g. for GCP: `mygcpproject`,

* `provider_type` - (Required) The cloud provider in which the account exists. Currently supported provider is `PROVIDER_GCP`.

* `enabled` - (Required) Whether or not to enable sysdig provisioning of resources on this cloud account.

* `feature` - (Optional) The name and configuration of each feature along with the respective components to enable on this cloud account.

* `component` - (Optional) The component configuration to enable on this cloud account. There can be multiple component blocks for a feature, one for each component to be enabled.

## Attributes Reference

No additional attributes are exported.
