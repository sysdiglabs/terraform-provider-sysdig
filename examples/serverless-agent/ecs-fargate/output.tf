output "workload_cluster_name" {
  value = aws_ecs_cluster.cluster.name
}

output "workload_cluster_arn" {
  value = aws_ecs_cluster.cluster.arn
}

output "service_arn" {
  value = aws_ecs_service.service.id
}

output "task_revision" {
  value = aws_ecs_task_definition.task_definition.revision
}
