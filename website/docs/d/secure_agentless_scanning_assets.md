---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_agentless_scanning_assets"
description: |-
  Retrieves information about the Sysdig Secure Agentless Scanning Assets
---

# Data Source: sysdig_secure_agentless_scanning_assets

Retrieves information about the Sysdig Secure Agentless Scanning Assets

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_agentless_scanning_assets" "assets" {}
```

## Argument Reference

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `aws.account_id` - AWS account sandbox in which Sysdig Agentless Scanning operates

* `azure.service_principal_id` - Azure service principal id for use with Sysdig Agentless Scanning

* `azure.tenant_id` - Azure tenant id in which Sysdig Agentless Scanning operates

* `backend.cloud_id` - Sysdig backend cloud identifier

* `backend.type` - Sysdig backend cloud type

* `gcp.worker_identity` - GCP worker indentity id

