---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_cloud_ingestion_assets"
description: |-
  Retrieves information about the Sysdig Secure Cloud Ingestion Assets
---

# Data Source: sysdig_secure_cloud_ingestion_assets

Retrieves information about the Sysdig Secure Cloud Ingestion Assets

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_cloud_ingestion_assets" "assets" {}
```

## Argument Reference

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `aws.eventBusARN` - AWS event bus from which Sysdig Cloud Ingestion operates

* `aws.eventBusARNGov` - AWS Gov event bus (if supported) from which Sysdig Cloud Ingestion operates

* `aws.sns_routing_key` - AWS CloudTrail SNS ingestion routing key

* `aws.sns_metadata` - AWS CloudTrail SNS ingestion metadata

* `gcp_routing_key` - GCP ingestion routing key

* `gcp_metadata` - GCP ingestion metadata
