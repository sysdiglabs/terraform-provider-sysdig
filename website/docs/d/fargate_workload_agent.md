---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_fargate_workload_agent"
description: |-
  Updates the fargate workload definition to add a [Sysdig Agent](https://docs.sysdig.com/en/docs/installation/serverless-agents/aws-fargate-serverless-agents/)
---

# Data Source: fargate_workload_agent

Updates the fargate workload definition to add a [Sysdig Agent](https://registry.terraform.io/modules/sysdiglabs/fargate-orchestrator-agent/aws/latest)

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

You can connect the Sysdig agent directly to a collector (using `collector_host` and `collector` port) or through an orchestrator (using `orchestrator_host` and `orchestrator_port`). For details about how to deploy an orchestrator check the [Sysdig Orchestrator module](https://registry.terraform.io/modules/sysdiglabs/fargate-orchestrator-agent/aws/latest).

## Example Usage

```terraform
data "sysdig_fargate_workload_agent" "task_definition" {
  container_definitions = "[]"

  image_auth_secret    = ""
  workload_agent_image = "quay.io/sysdig/workload-agent:latest"

  sysdig_access_key = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  collector_host    = "collector-static.sysdigcloud.com"
  collector_port    = 6443
}
```

## Argument Reference

* `container_definitions` - (Required) The input Fargate container definitions to instrument with the Sysdig workload agent.
* `sysdig_access_key` - (Required) The Sysdig Access Key (Agent token).
* `workload_agent_image` - (Required) The Sysdig workload agent image.
* `image_auth_secret` - (Optional) The registry authentication secret.
* `orchestrator_host` - (Optional) The orchestrator host to connect to.
* `orchestrator_port` - (Optional) The orchestrator port to connect to.
* `collector_host` - (Optional) The collector host to connect to.
* `collector_port` - (Optional) The collector port to connect to.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `output_container_definitions` - The updated container definitions instrumented with the Sysdig workload agent.
