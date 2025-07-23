---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_fargate_workload_agent"
description: |-
  Updates the fargate workload definition to add a Sysdig Agent
---

# Data Source: sysdig_fargate_workload_agent

Updates the ECS Fargate Container Definitions to add a [Sysdig Workload Agent](https://docs.sysdig.com/en/docs/sysdig-secure/install-agent-components/linux-on-serverless/ecs-fargate/)

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

The Sysdig Workload Agent will need to connect to the Sysdig Collector. Find your region's collector endpoint here: https://docs.sysdig.com/en/docs/administration/saas-regions-and-ip-ranges/.

## Example Usage

```terraform
data "sysdig_fargate_workload_agent" "instrumented_containers" {
  container_definitions = "[]"

  workload_agent_image = "quay.io/sysdig/workload-agent:latest"

  collector_host    = var.collector_host
  collector_port    = var.collector_port
  sysdig_access_key = var.sysdig_access_key
}
```

## Argument Reference

* `container_definitions` - (Required) The input Fargate container definitions to instrument with the Sysdig workload agent.
* `workload_agent_image` - (Required) The Sysdig workload agent image.
* `collector_host` - (Required) The Sysdig Collector host to connect to.
* `collector_port` - (Required) The Sysdig Collector port.
* `sysdig_access_key` - (Required) The Sysdig Agent access key, available in the Sysdig Secure UI.
* `image_auth_secret` - (Optional) The registry authentication secret.
* `log_configuration` - (Optional) Configuration for the awslogs driver on the instrumentation container. All three values must be specified if instrumentation logging is desired:
  * `group` - The name of the existing log group for instrumentation logs
  * `stream_prefix` - Prefix for the instrumentation log stream
  * `region` - The AWS region where the target log group resides
* `sysdig_logging` - (Optional) The instrumentation logging level: `trace`, `debug`, `info`, `warning`, `error`, `silent`.
* `ignore_containers` - (Optional) A list of containers in this data source that should not be instrumented.
* `bare_pdig_on_containers` - (Optional) A list of containers in this data source to be instrumented with bare pdig.
* `priority` - (Optional) The priority mode for the workload agent. Can be `availability` (by default) or `security`.
* `instrumentation_essential` - (Optional) `false` by default. If `true` the instrumentation container will be marked as essential.
* `instrumentation_cpu` - (Optional) The number of CPU units for the instrumentation container.
* `instrumentation_memory_limit` - (Optional) The maximum amount (in MiB) of memory for the instrumentation container.
* `instrumentation_memory_reservation` - (Optional) The minimum amount (in MiB) of memory reserved for the instrumentation container.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `output_container_definitions` - The updated container definitions instrumented with the Sysdig workload agent.
