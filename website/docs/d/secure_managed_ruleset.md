---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_managed_ruleset"
description: |-
  Retrieves a Sysdig Secure Managed Ruleset.
---

# Data Source: sysdig_secure_managed_ruleset

Retrieves the information of an existing Sysdig Secure Managed Ruleset.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_managed_ruleset" "example" {
  name                 = "Sysdig Runtime Threat Detection - Managed Ruleset"
  type                 = "falco"
}
```

## Argument Reference

* `name` - (Required) The name of the Secure managed ruleset.

* `type` - (Optional) Specifies the type of the runtime policy. Must be one of: `falco`, `list_matching`, `k8s_audit`,
  `aws_cloudtrail`, `gcp_auditlog`, `azure_platformlogs`. By default it is `falco`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id for the managed policy.

* `description` - The description for the managed policy.

* `severity` -  The severity of Secure policy. The accepted values
    are: 0, 1, 2, 3 (High), 4, 5 (Medium), 6 (Low) and 7 (Info).

* `enabled` - Whether the policy is enabled or not.

* `runbook` - Customer provided url that provides a runbook for a given policy.

* `scope` - The application scope for the policy.

* `rules` - An array of rules with the properties `name` and `enabled` to identify the rule name and whether it is enabled.

* `notification_channels` - IDs of the notification channels to send alerts to
    when the policy is fired.

### Actions block

The actions block is optional and supports:

* `container` - (Optional) The action applied to container when this Policy is
    triggered. Can be *stop*, *pause* or *kill*. If this is not specified,
    no action will be applied at the container level.

* `capture` - (Optional) Captures with Sysdig the stream of system calls:
    * `seconds_before_event` - (Required) Captures the system calls during the
    amount of seconds before the policy was triggered.
    * `seconds_after_event` - (Required) Captures the system calls for the amount
    of seconds after the policy was triggered.
    * `name` - (Required) The name of the capture file
    * `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name=cat`
    * `bucket_name` - (Optional) Custom bucket to store capture in, 
    bucket should be onboarded in Integrations > S3 Capture Storage. Default is to use Sysdig Secure Storage 
    * `folder` - (Optional) Name of folder to store capture inside the bucket. 
    By default we will store the capture file at the root of the bucket
