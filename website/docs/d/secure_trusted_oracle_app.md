---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_trusted_oracle_app"
description: |-
  Retrieves information about the Sysdig Secure Trusted Oracle App
---

# Data Source: sysdig_secure_trusted_oracle_app

Retrieves information about the Sysdig Secure Trusted Oracle App

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_trusted_oracle_app" "onboarding" {
	name = "onboarding"
}
```

## Argument Reference

* `name` - (Required) Sysdig's Oracle App name. Currently supported applications are `config_posture` and `onboarding`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tenancy_ocid` - The application's associated tenancy OCI identifer.

* `group_ocid` - The application's associated usergroup OCI identifier.

* `user_ocid` - The application's associated user OCI identifier.

