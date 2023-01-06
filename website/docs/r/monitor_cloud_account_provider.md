---
subcategory: "Sysdig Monitor"
layout: "sysdig"
page_title: "Sysdig: sysdig_monitor_cloud_account_provider"
description: |-
Creates a Sysdig Monitor cloud provider integration
---

# Resource: sysdig_monitor_cloud_account_provider

Creates a Sysdig Monitor cloud provider integration used to monitor GCP cloud accounts.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_monitor_cloud_account_provider" "sample" {
  platform = "GCP"
  integration_type = "API"
  account_id = "A60E650B-B24F-4934-867A-45D5F0DB814E"
}
```

## Argument Reference

### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `platform` - (Required) Cloud platform that will be monitored. Only `GCP` is currently supported.
* `integration_type` - (Required) Type of cloud integration. Only `API` is currently supported.
* `account_id` - (Required) The GCP project id for the project that will be monitored.
* `additional_options` - (Optional) The private key generated when creating a new GCP service account key. Must be in JSON format and base64 encoded.