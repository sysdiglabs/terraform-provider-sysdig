output "orchestrator_cluster_name" {
  value = data.aws_ecs_cluster.fargate-orchestrator.cluster_name
}

output "orchestrator_cluster_arn" {
  value = data.aws_ecs_cluster.fargate-orchestrator.arn
}

output "orchestrator_service_arn" {
  value = data.aws_ecs_service.orchestrator-service.arn
}
