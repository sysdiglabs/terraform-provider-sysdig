---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_trusted_cloud_identity"
description: |-
  Retrieves information about the Sysdig Secure Trusted Cloud Identity
---

# Data Source: sysdig_secure_trusted_cloud_identity

Retrieves information about the Sysdig Secure Trusted Cloud Identity

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example Usage

```hcl
data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {
	cloud_provider = "aws"
}
```

## Argument Reference

* `cloud_provider` - (Required) The cloud provider in which the trusted identity will be used. Currently supported providers are `aws`, `gcp` and `azure` 


## Attributes Reference

* `identity` - Sysdig's identity (User/Role/etc) that should be used to create a trust relationship allowing Sysdig access to your cloud account.

* `aws_account_id` - If `identity` is an AWS ARN, this attribute contains the AWS Account ID to which the ARN belongs, otherwise it contains the empty string.

* `aws_role_name` - If `identity` is a AWS IAM Role ARN, this attribute contains the name of the role, otherwise it contains the empty string.
