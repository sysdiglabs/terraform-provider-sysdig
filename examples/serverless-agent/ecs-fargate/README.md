# Workload with Serverless Workload Agent
This example deploys a cluster running a workload application secured by the Serverless Workload Agent.
The workload used is [falcosecurity/event-generator](https://github.com/falcosecurity/event-generator), which produces synthetic suspicious actions that will trigger Sysdig managed policies. 


## Prerequisites
The following prerequisites are required to deploy this sample:
- `VPC ID`, the ID of an already existing VPC
- `Subnet ID`, the ID of an already existing Subnet within the VPC above


## Usage
```
$ terraform init
$ terraform apply
```


## Components
| **Component**      | **Name**                   | **Description**                                                                                                                                     |
|--------------------|----------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| ECS Cluster        | `<prefix>-cluster`         | The cluster containing the ECS Service below.                                                                                                       |
| ECS Service        | `<prefix>-service`         | The service running the TaskDefinition below.                                                                                                       |
| ECS TaskDefinition | `<prefix>-task-definition` | The task definition including a workload container being secured by the Serverless Agent.                                                           |
| ECS SecurityGroup  | `<prefix>-security-group`  | The security group ensuring connectivity to the Serverless Agent. This security group has no restrictions applied and is intended for testing only. |


## Files
| **File**       | **Description**                                                                                          |
|----------------|----------------------------------------------------------------------------------------------------------|
| `output.tf`    | Contains the reference to the cluster, service, and task revision being deployed.                        |
| `providers.tf` | Contains the configuration parameters for the providers.                                                 |
| `resources.tf` | Contains the resources to deploy, including the task definition being secured with the Serverless Agent. |
| `variables.tf` | Contains the configuration parameters for AWS and the Serverless Agent.                                  |
| `versions.tf`  | Defines the version of the providers.                                                                    |
