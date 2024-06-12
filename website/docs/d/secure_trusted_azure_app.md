---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_trusted_azure_app"
description: |-
  Retrieves information about the Sysdig Secure Trusted Azure App
---

# Data Source: sysdig_secure_trusted_azure_app

Retrieves information about the Sysdig Secure Trusted Azure App

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_trusted_azure_app" "onboarding" {
	name = "onboarding"
}
```

## Argument Reference

* `name` - (Required) Sysdig's Azure App name urrently supported applications are `config_posture`, `onboarding` and `threat_detection` 


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tenant_id` - The application's associated tenant identifer

* `application_id` - The application's identifier

* `service_principal_id` - The application's associated service principal identifier.

