# AWS configuration
variable "prefix" {
  description = "All resources created by Terraform have this prefix prepended to them"
}

variable "profile" {
  description = "AWS profile name"
  type        = string
}

variable "region" {
  description = "AWS Region for deployment"
  default     = "us-east-1"
}

variable "subnet_1" {
  description = "Subnet-1 Id"
}

variable "subnet_2" {
  description = "Subnet-2 Id"
}

variable "vpc_id" {
  description = "VPC Id"
}

variable "tags" {
  type        = map(string)
  description = "Tags to assign to resources in module"
  default     = {}
}

variable "replicas" {
  description = "Number of workload replicas to run"
  default     = 2
}

# Serverless Agent Configuration
variable "access_key" {
  description = "Sysdig Agent access key"
}

variable "agent_workload_image" {
  description = "Workload agent container image"
  default     = "quay.io/sysdig/workload-agent:latest"
}

variable "collector_host" {
  description = "Collector Host"
}

variable "collector_port" {
  description = "Collector Port"
  default     = 6443
}
