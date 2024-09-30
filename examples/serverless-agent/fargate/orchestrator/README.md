# Serverless Orchestrator Agent

This example deploys an AWS ECS Fargate cluster to run the Serverless Orchestrator Agent. This Agent acts as a proxy between the Collector and many Serverless Workload Agents.

## Prerequisites

The following AWS prerequisites are required to deploy this cluster:
- VPC
- 2 subnets

## Components

The cluster will be called `<prefix>-cluster` and will deploy the following:
- 1 Service (called `OrchestratorAgent`)
  - 1 Task (with the latest version of the Serverless Orchestrator Agent)
- Network Load balancer
- Cloudwatch log group
- Security group

## Layout
| **File** | **Purpose** |
| --- | --- |
| `main.tf` | AWS provider configuration |
| `orchestrator.tf` | Orchestrator cluster definition |
| `output.tf` | Defines the output variables |
| `variables.tf` | AWS and Agent configuration |
| `versions.tf` | Defines TF provider versions |
