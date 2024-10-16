---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_trusted_cloud_regulation_assets"
description: |-
  Retrieves information about the Sysdig Secure Trusted Cloud Regulation Assets
---

# Data Source: sysdig_secure_trusted_cloud_regulation_assets

Retrieves information about the Sysdig Secure Trusted Cloud Regulation Assets

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_trusted_cloud_regulation_assets" "trusted_identity_gov" {
	cloud_provider = "aws"
}
```

## Argument Reference

* `cloud_provider` - (Required) The cloud provider in which the trusted identity for regulatory workloads will be used. Currently supported providers are `aws` 


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gov_identity` - Sysdig's identity for regulatory workloads (User/Role/etc) that should be used to create a trust relationship allowing Sysdig access to your regulated cloud account.

* `aws_gov_account_id` - If `gov_identity` is an AWS GOV ARN, this attribute contains the AWS GOV Account ID to which the ARN belongs, otherwise it contains the empty string. `cloud_provider` must be equal to `aws`.

* `aws_gov_role_name` - If `gov_identity` is a AWS GOV IAM Role ARN, this attribute contains the name of the GOV role, otherwise it contains the empty string. `cloud_provider` must be equal to `aws`.

