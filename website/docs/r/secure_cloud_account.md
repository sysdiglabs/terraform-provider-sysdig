---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_cloud_account"
description: |-
  Creates a Sysdig Secure Cloud Account.
---

# Resource: sysdig_secure_cloud_account

Creates a Sysdig Secure Cloud Account.

~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_cloud_account" "sample" {
  account_id          = "123456789012"
  cloud_provider      = "aws"
  alias               = "prod"
  role_enabled        = "false"
  role_name           = "CustomRoleName"
}
```

## Argument Reference

* `account_id` - (Required) The unique identifier of the cloud account. e.g. for AWS: `123456789012`,

* `cloud_provider` - (Required) The cloud provider in which the account exists. Currently supported providers are `aws`, `gcp` and `azure`

* `alias` - (Optional) A human friendly alias for `account_id`.

* `role_enabled` - (Optional) Whether or not a role is provisioned withing this account, that Sysdig has permission to AssumeRole in order to run Benchmarks. Default: `false`.

* `role_name` - (Optional) The name of the role Sysdig will have permission to AssumeRole if `role_enaled` is set to `true`. Default: `SysdigCloudBench`.

## Attributes Reference

No additional attributes are exported.

## Import

Secure Cloud Accounts can be imported using the `account_id`, e.g.

```
$ terraform import sysdig_secure_cloud_account.sample 123456789012
```
