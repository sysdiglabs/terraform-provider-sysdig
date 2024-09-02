module "fargate-orchestrator-agent" {
  source  = "sysdiglabs/fargate-orchestrator-agent/aws"
  version = "0.5.0"

  vpc_id  = var.vpc_id
  subnets = [var.subnet_1, var.subnet_2]

  access_key = var.access_key

  collector_host = var.collector_host
  collector_port = var.collector_port

  name        = var.prefix
  agent_image = var.agent_orchestrator_image

  # True if the VPC uses an InternetGateway, false otherwise
  assign_public_ip = true

  tags = var.tags
}


data "aws_ecs_cluster" "fargate-orchestrator" {
  depends_on = [
    module.fargate-orchestrator-agent
  ]
  cluster_name = "${var.prefix}-cluster"
}

data "aws_ecs_service" "orchestrator-service" {
  depends_on = [
    module.fargate-orchestrator-agent
  ]
  service_name = "OrchestratorAgent"
  cluster_arn  = data.aws_ecs_cluster.fargate-orchestrator.arn
}
