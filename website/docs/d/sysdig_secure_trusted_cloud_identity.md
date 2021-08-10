---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_trusted_cloud_identity"
sidebar_current: "docs-sysdig-secure-trusted-cloud-identity-ds"
description: |-
  Retrieves information about the Sysdig Secure Trusted Cloud Identity
---

# sysdig\_secure\_trusted_cloud_identity

Retrieves information about the Sysdig Secure Trusted Cloud Identity

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {
	cloud_provider = "aws"
}
```

## Argument Reference

* `cloud_provider` - (Required) The cloud provider in which the account exists. Currently supported providers are `aws`, `gcp` and `azure` 


## Attributes Reference

* `identity` - Sysdig's identity (User/Role/etc) that should be used to create a trust relationship allowing Sysdig access to your cloud account.
