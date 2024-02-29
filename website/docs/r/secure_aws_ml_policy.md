---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_aws_ml_policy"
description: |-
  Retrieves a Sysdig Secure AWS ML Policy.
---

# Resource: sysdig_secure_aws_ml_policy

Retrieves the information of an existing Sysdig Secure ML Policy.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_aws_ml_policy" "policy" {
  name        = "Test ML Policy 1"
  description = "Test ML Policy Description"
  enabled     = true
  severity    = 4

  rule {
    description = "Test ML Rule Description"

    anomalous_console_login {
      enabled   = true
      threshold = 1
      severity  = 1
    }
}
```

## Argument Reference

* `name` - (Required) The name of the Secure managed policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id for the policy.

* `description` - The description for the managed policy.

* `severity` -  The severity of Secure policy. The accepted values
    are: 0, 1, 2, 3 (High), 4, 5 (Medium), 6 (Low) and 7 (Info).

* `enabled` - Whether the policy is enabled or not.

* `runbook` - Customer provided url that provides a runbook for a given policy.

* `scope` - The application scope for the policy.

* `notification_channels` - IDs of the notification channels to send alerts to
    when the policy is fired.

### `rule` block

The rule block is required and supports:

* `description` - (Required) Rule description.
* `anomalous_console_login` - (Required) This attribute allows you to activate anomaly detection for console logins and adjust its settings.
    * `threshold` - (Required) Trigger at or above confidence level.

