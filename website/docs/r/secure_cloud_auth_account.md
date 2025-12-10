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

### Basic Usage

```terraform
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "mygcpproject"
  provider_type = "PROVIDER_GCP"
  enabled       = true
}
```

### AWS Account with Response Actions

```terraform
resource "sysdig_secure_cloud_auth_account" "aws_response_actions" {
  provider_id   = "123456789012"
  provider_type = "PROVIDER_AWS"
  enabled       = true

  feature {
    secure_response_actions {
      enabled    = true
      components = [
        "COMPONENT_CLOUD_RESPONDER/cloud-responder",
        "COMPONENT_CLOUD_RESPONDER_ROLES/cloud-responder"
      ]
    }
  }

  component {
    type                     = "COMPONENT_CLOUD_RESPONDER"
    instance                 = "cloud-responder"
    cloud_responder_metadata = jsonencode({
      aws = {
        responder_lambdas = {
          lambda_names       = ["sysdig-responder-lambda-1", "sysdig-responder-lambda-2"]
          regions            = ["us-east-1", "us-west-2"]
          delegate_role_name = "sysdig-delegate-role"
        }
      }
    })
  }

  component {
    type                           = "COMPONENT_CLOUD_RESPONDER_ROLES"
    instance                       = "cloud-responder"
    cloud_responder_roles_metadata = jsonencode({
      roles = [
        { aws = { role_name = "sysdig-responder-role-1" } },
        { aws = { role_name = "sysdig-responder-role-2" } }
      ]
    })
  }
}
```

## Argument Reference

* `provider_id` - (Required) The unique identifier of the cloud account. e.g. for GCP: `mygcpproject`.

* `provider_type` - (Required) The cloud provider in which the account exists. Currently supported provider is `PROVIDER_GCP`.

* `enabled` - (Required) Whether or not to enable sysdig provisioning of resources on this cloud account.

* `feature` - (Optional) The name and configuration of each feature along with the respective components to enable on this cloud account.

* `component` - (Optional) The component configuration to enable on this cloud account. There can be multiple component blocks for a feature, one for each component to be enabled.

* `cloud_responder_metadata` - (Optional) Configuration metadata for the Cloud Responder component (type `COMPONENT_CLOUD_RESPONDER`). Used with the Response Actions feature to specify Lambda functions and IAM roles for automated response capabilities.

* `cloud_responder_roles_metadata` - (Optional) Configuration metadata for the Cloud Responder Roles component (type `COMPONENT_CLOUD_RESPONDER_ROLES`). Defines the IAM roles that the cloud responder can assume when executing response actions.

* `provider_partition` - (Optional - AWS installs only) The type of Partition of the Provider for cloud account. Currently supported options are `PROVIDER_PARTITION_UNSPECIFIED` and `PROVIDER_PARTITION_AWS_GOVCLOUD`.


-> **Note:** Please refer to Sysdig Secure API Documentation for the Cloud Accounts API for providing `feature` & `component`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Computed) The ID of the cloud account.

* `organization_id` - (Computed) The ID of the organization, if the cloud account is part of any organization.