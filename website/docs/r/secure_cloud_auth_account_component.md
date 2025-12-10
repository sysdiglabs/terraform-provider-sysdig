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

### GCP Service Principal Component

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

### AWS Cloud Responder Component

```terraform
resource "sysdig_secure_cloud_auth_account" "aws_account" {
  provider_id   = "123456789012"
  provider_type = "PROVIDER_AWS"
  enabled       = true
  lifecycle {
    ignore_changes = [
      component,
      feature
    ]
  }
}

resource "sysdig_secure_cloud_auth_account_component" "cloud_responder" {
  account_id               = sysdig_secure_cloud_auth_account.aws_account.id
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
```

### AWS Cloud Responder Roles Component

```terraform
resource "sysdig_secure_cloud_auth_account" "aws_account" {
  provider_id   = "123456789012"
  provider_type = "PROVIDER_AWS"
  enabled       = true
  lifecycle {
    ignore_changes = [
      component,
      feature
    ]
  }
}

resource "sysdig_secure_cloud_auth_account_component" "cloud_responder_roles" {
  account_id                     = sysdig_secure_cloud_auth_account.aws_account.id
  type                           = "COMPONENT_CLOUD_RESPONDER_ROLES"
  instance                       = "cloud-responder"
  cloud_responder_roles_metadata = jsonencode({
    roles = [
      { aws = { role_name = "sysdig-responder-role-1" } },
      { aws = { role_name = "sysdig-responder-role-2" } },
      { aws = { role_name = "sysdig-responder-role-3" } }
    ]
  })
}
```

## Argument Reference

* `account_id` - (Required) Cloud Account created using resource sysdig_secure_cloud_auth_account.

* `type` - (Required) The type of component to be created. e.g. `COMPONENT_SERVICE_PRINCIPAL`.

* `instance` - (Required) The component instance to be created, identified by a specific string. e.g. `secure-posture`, `secure-runtime`, `cloud-responder`, etc.

* `<component>_metadata` - (Optional) Based on the component type created, this is the metadata information passed to enable the component on the account.

* `cloud_responder_metadata` - (Optional) Metadata for `COMPONENT_CLOUD_RESPONDER` type. Configures the Lambda functions and IAM roles for automated response actions. Required fields:
  * `aws.responder_lambdas.lambda_names` - List of Lambda function names to use for response actions
  * `aws.responder_lambdas.regions` - List of AWS regions where the responder is deployed
  * `aws.responder_lambdas.delegate_role_name` - IAM role name that the responder assumes

* `cloud_responder_roles_metadata` - (Optional) Metadata for `COMPONENT_CLOUD_RESPONDER_ROLES` type. Defines the IAM roles that can be assumed for response actions. Required fields:
  * `roles` - Array of role objects, each containing provider-specific role name (e.g., `aws.role_name` for AWS roles)

-> **Note:** Please refer to Sysdig Secure API Documentation for the Cloud Accounts API for metadata types for `component`.

-> **Note:** Since creation of component resource updates the account resource in the backend, in these configurations we indicate to Terraform to ignore `component` & `feature` attributes when planning updates to the remote account resource object.

## Attributes Reference

No additional attributes are exported.