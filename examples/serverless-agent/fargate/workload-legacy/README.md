# Workload with Serverless Workload Agent

This example deploys a cluster with a workload and the Serverless Workload Agent as a sidecar to secure the workload.

The Workload Agent will use an Orchestrator Agent as a proxy to the Sysdig Collector.

## Prerequisites

The following prerequisites are required to deploy this cluster:
- Orchestrator Agent deployed
- VPC
- 2 subnets

## Components

The cluster will be called `<prefix>-instrumented-workload` and will deploy the following:
- 1 Service (called `<prefix-instrumented-service`)
  - 1 Task with 2 replicas, each running:
    - 1 container named `event-gen-1` running `falcosecurity/event-generator`
    - 1 container named `event-gen-2` also running `falcosecurity/event-generator`
    - 1 container named `SysdigInstrumentation` running the latest Workload Agent which will secure both workload containers

## Layout
| **File** | **Purpose** |
| --- | --- |
| `instrumented_load.tf` | Workload definition. By default it instruments `falcosecurity/event-generator` |
| `main.tf` | AWS provider configuration |
| `output.tf` | Defines the output variables |
| `variables.tf` | AWS and Agent configuration |
| `versions.tf` | Defines TF provider versions |
