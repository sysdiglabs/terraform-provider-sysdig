---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_benchmark_task"
description: |-
  Creates a Sysdig Secure Benchmark Task.
---

# Resource: sysdig\_secure\_benchmark_task

Creates a Sysdig Secure Benchmark Task.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_secure_benchmark_task" "sample" {
  name     = "My Benchmark Task"
  schedule = "0 6 * * *"
  schema   = "aws_foundations_bench-1.3.0"
  scope    = "aws.accountId = \"123456789012\" and aws.region = \"us-west-2\""
  enabled  = "true"
}
```

## Argument Reference

* `name` - (Required) The unique identifier of the cloud account. e.g. for AWS: `123456789012`, 

* `schedule` - (Required) The schedule (as a cron expression: [Minute Hour Day DayOfWeek DayOfMonth]) on which this task should be run. The schedule may not be more frequent than once per day.

* `schema` - (Required) The identifier of the benchmark schema of which to run. Possible values are: `aws_foundations_bench-1.3.0`, `gcp_foundations_bench-1.2.0`, `azure_foundations_bench-1.3.0`.

* `scope` - (Required) The Sysdig scope expression on which to run this benchmark: e.g. `aws.accountId = \"123456789012\" and aws.region = \"us-west-2\"`. The labels available are `aws.accountId`, `aws.region`, `gcp.projectId` and `gcp.region`. Only the `=` and `and` operators are supported.

* `enabled` - (Optional) Whether or not this task should be enabled. Default: `true`.

## Import

Secure Benchmark Tasks can be imported using the `id`, e.g.

```
$ terraform import sysdig_secure_benchmark_task.sample 1
```
