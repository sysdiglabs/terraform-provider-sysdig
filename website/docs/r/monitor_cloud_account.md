---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_cloud_account"
description: |- 
  Creates a Sysdig Monitor Cloud Account 
---

# Resource: sysdig_monitor_cloud_account

Creates a Sysdig Monitor Cloud Account for monitoring cloud resources.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
// GCP example
resource "sysdig_monitor_cloud_account" "sample" {
  cloud_provider = "GCP"
  integration_type = "API"
  account_id = "gcp_project_id"
}

// AWS example with role delegation
resource "sysdig_monitor_cloud_account" "sample" {
  cloud_provider = "AWS"
  integration_type = "Metrics Streams"
  account_id = "123412341234"
  role_name = "SysdigTestRole"
}

// AWS example with secret key
resource "sysdig_monitor_cloud_account" "sample" {
  cloud_provider = "AWS"
  integration_type = "Metrics Streams"
  account_id = "123412341234"
  secret_key = "Xxx5XX2xXx/Xxxx+xxXxXXxXxXxxXXxxxXXxXxXx"
  access_key_id = "XXXXX33XXXX3XX3XXX7X"
}
```

## Argument Reference

* `cloud_provider` - (Required) Cloud platform that will be monitored. Only `GCP` and `AWS` are currently supported.
* `integration_type` - (Required) Type of cloud integration. Only `API` and `Metrics Streams` are currently supported (`Metrics Streams` only for `AWS`).
* `account_id` - (Required) The GCP project id for the project that will be monitored and Account ID itself for AWS.
* `role_name` - (Optional) The role name tha will take the permissions over some resources in AWS
* `secret_key` - (Optional) The the secret key for a AWS connection
* `access_key_id` - (Optional) The ID for the access key that has the permissions into the Cloud Account
* `additional_options` - (Optional) The private key generated when creating a new GCP service account key. Must be in JSON format and base64 encoded.

## Attributes Reference

No additional attributes are exported.
