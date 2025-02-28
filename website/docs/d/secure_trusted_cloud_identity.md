---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_trusted_cloud_identity"
description: |-
  Retrieves information about the Sysdig Secure Trusted Cloud Identity
---

# Data Source: sysdig_secure_trusted_cloud_identity

Retrieves information about the Sysdig Secure Trusted Cloud Identity

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_trusted_cloud_identity" "trusted_identity" {
	cloud_provider = "aws"
}
```

## Argument Reference

* `cloud_provider` - (Required) The cloud provider in which the trusted identity will be used. Currently supported providers are `aws`, `gcp` and `azure` 


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `identity` - Sysdig's identity (User/Role/etc) that should be used to create a trust relationship allowing Sysdig access to your cloud account.

* `aws_account_id` - If `identity` is an AWS ARN, this attribute contains the AWS Account ID to which the ARN belongs, otherwise it contains the empty string. `cloud_provider` must be equal to `aws` or `gcp`.

* `aws_role_name` - If `identity` is a AWS IAM Role ARN, this attribute contains the name of the role, otherwise it contains the empty string. `cloud_provider` must be equal to `aws` or `gcp`.

* `azure_tenant_id` - If `identity` contains credentials for an Azure Service Principal, this attribute contains the service principal's Tenant ID. `cloud_provider` must be equal to `azure`.

* `azure_service_principal_id` - If `identity` contains credentials for an Azure Service Principal, this attribute contains the service principal's ID. `cloud_provider` must be equal to `azure`.

* `gov_identity` - Sysdig's identity for regulatory workloads (User/Role/etc) that should be used to create a trust relationship allowing Sysdig access to your regulated cloud account. Currently supported on `aws`.

* `aws_gov_account_id` - If `gov_identity` is an AWS GOV IAM Role ARN, this attribute contains the AWS GOV Account ID to which the ARN belongs, otherwise it contains the empty string. Currently supported on `aws`.

* `aws_gov_role_name` - If `gov_identity` is a AWS GOV IAM Role ARN, this attribute contains the name of the GOV role, otherwise it contains the empty string. Currently supported on `aws`.
