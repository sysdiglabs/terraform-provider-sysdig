---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_okta_ml_policy"
description: |-
  Retrieves a Sysdig Secure Okta ML Policy.
---

# Data Source: sysdig_secure_okta_ml_policy

Retrieves information about an existing Sysdig Secure Okta ML Policy.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_okta_ml_policy" "policy" {
  name = "My Okta ML Policy"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure Okta ML policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id for the policy.

* `type` - The type of the policy. Always set to "okta_machine_learning".

* `description` - The description for the policy.

* `severity` -  The severity of Secure policy. The accepted values
    are: 0, 1, 2, 3 (High), 4, 5 (Medium), 6 (Low) and 7 (Info).

* `enabled` - Whether the policy is enabled or not.

* `runbook` - Customer provided url that provides a runbook for a given policy.

* `scope` - The application scope for the policy.

* `notification_channels` - IDs of the notification channels to send alerts to
    when the policy is fired.

### `rule` block

The rule block contains:

* `id` - The ID of the rule.

* `name` - The name of the rule.

* `description` - Rule description.

* `tags` - Tags associated with the rule.

* `version` - The version of the rule.

* `anomalous_console_login` - Anomaly detection settings for logins.
    * `enabled` - Whether anomaly detection is enabled.
    * `threshold` - Confidence level threshold for triggering alerts. Valid values are: 1 (Default), 2 (High), 3 (Higher).
