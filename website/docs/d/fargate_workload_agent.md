---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_fargate_workload_agent"
description: |-
  Updates the fargate workload definition to add a Sysdig Agent
---

# Data Source: fargate_workload_agent

Updates the fargate workload definition to add a [Sysdig Agent](https://docs.sysdig.com/en/docs/installation/serverless-agents/aws-fargate-serverless-agents/)

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

You'll need to connect the Sysdig Agent to the Sysdig backend through an orchestrator. For details about how to deploy an orchestrator check the [Sysdig Orchestrator module](https://registry.terraform.io/modules/sysdiglabs/fargate-orchestrator-agent/aws/latest).

## Example Usage

```terraform
data "sysdig_fargate_workload_agent" "instrumented_containers" {
  container_definitions = "[]"

  image_auth_secret    = ""
  workload_agent_image = "quay.io/sysdig/workload-agent:latest"

  sysdig_access_key = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  orchestrator_host    = module.fargate-orchestrator-agent.orchestrator_host
  orchestrator_port    = module.fargate-orchestrator-agent.orchestrator_port
}
```

## Argument Reference

* `container_definitions` - (Required) The input Fargate container definitions to instrument with the Sysdig workload agent.
* `sysdig_access_key` - (Required) The Sysdig Access Key (Agent token).
* `orchestrator_host` - (Required) The orchestrator host to connect to.
* `orchestrator_port` - (Required) The orchestrator port to connect to.
* `workload_agent_image` - (Required) The Sysdig workload agent image.
* `image_auth_secret` - (Optional) The registry authentication secret.
* `log_configuration` - (Optional) Configuration for the awslogs driver on the instrumentation container. All three values must be specified if instrumentation logging is desired:
  * `group` - The name of the existing log group for instrumentation logs
  * `stream_prefix` - Prefix for the instrumentation log stream
  * `region` - The AWS region where the target log group resides
* `sysdig_logging` - (Optional) The instrumentation logging level: `trace`, `debug`, `info`, `warning`, `error`, `silent`.
* `ignore_containers` - (Optional) A list of containers in this data source that should not be instrumented.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `output_container_definitions` - The updated container definitions instrumented with the Sysdig workload agent.
