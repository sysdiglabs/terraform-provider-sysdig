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

# Serverless Agent Configuration
variable "access_key" {
  description = "Sysdig Agent access key"
}

variable "agent_orchestrator_image" {
  description = "Orchestrator Agent image to use"
  default     = "quay.io/sysdig/orchestrator-agent:latest"
}

variable "collector_host" {
  description = "Collector host where agent will send the data"
}

variable "collector_port" {
  description = "Collector port where agent will send the data"
  default     = "6443"
}
