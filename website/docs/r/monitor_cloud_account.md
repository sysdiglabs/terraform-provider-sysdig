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
resource "sysdig_monitor_cloud_account" "sample" {
  cloud_provider = "GCP"
  integration_type = "API"
  account_id = "gcp_project_id"
}
```

## Argument Reference

### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `cloud_provider` - (Required) Cloud platform that will be monitored. Only `GCP` is currently supported.
* `integration_type` - (Required) Type of cloud integration. Only `API` is currently supported.
* `account_id` - (Required) The GCP project id for the project that will be monitored.
* `additional_options` - (Optional) The private key generated when creating a new GCP service account key. Must be in JSON format and base64 encoded.

## Attributes Reference

No additional attributes are exported.
