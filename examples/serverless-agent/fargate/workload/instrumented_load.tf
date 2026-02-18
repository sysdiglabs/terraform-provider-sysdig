data "sysdig_fargate_workload_agent" "containers_instrumented" {
  container_definitions = jsonencode([
    {
      "name" : "event-gen-1",
      "image" : "falcosecurity/event-generator",
      "command" : ["run", "syscall", "--all", "--loop"],
      "logConfiguration" : {
        "logDriver" : "awslogs",
        "options" : {
          "awslogs-group" : aws_cloudwatch_log_group.instrumented_logs.name,
          "awslogs-region" : var.region,
          "awslogs-stream-prefix" : "task"
        },
      }
    },
    {
      "name" : "event-gen-2",
      "image" : "falcosecurity/event-generator",
      "command" : ["run", "syscall", "--all", "--loop"],
      "logConfiguration" : {
        "logDriver" : "awslogs",
        "options" : {
          "awslogs-group" : aws_cloudwatch_log_group.instrumented_logs.name,
          "awslogs-region" : var.region,
          "awslogs-stream-prefix" : "task"
        },
      }
    }
  ])

  workload_agent_image = var.agent_workload_image

  sysdig_access_key = var.access_key
  collector_host    = var.collector_host
  collector_port    = var.collector_port

  log_configuration {
    group         = aws_cloudwatch_log_group.instrumented_logs.name
    stream_prefix = "instrumentation"
    region        = var.region
  }
}

resource "aws_ecs_task_definition" "task_definition" {
  family             = "${var.prefix}-instrumented-task-definition"
  task_role_arn      = aws_iam_role.task_role.arn
  execution_role_arn = aws_iam_role.execution_role.arn

  cpu                      = "256"
  memory                   = "512"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  pid_mode                 = "task"

  container_definitions = data.sysdig_fargate_workload_agent.containers_instrumented.output_container_definitions
}


resource "aws_ecs_cluster" "cluster" {
  name = "${var.prefix}-instrumented-workload"
}

resource "aws_cloudwatch_log_group" "instrumented_logs" {
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "execution_role" {
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json

  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"]
}

resource "aws_iam_role" "task_role" {
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json

  inline_policy {
    name   = "root"
    policy = data.aws_iam_policy_document.task_policy.json
  }
}

data "aws_iam_policy_document" "task_policy" {
  statement {
    actions = [
      "ecr:GetAuthorizationToken",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["*"]
  }
}

resource "aws_ecs_service" "service" {
  name = "${var.prefix}-instrumented-service"

  cluster          = aws_ecs_cluster.cluster.id
  task_definition  = aws_ecs_task_definition.task_definition.arn
  desired_count    = var.replicas
  launch_type      = "FARGATE"
  platform_version = "1.4.0"

  network_configuration {
    subnets          = [var.subnet_1, var.subnet_2]
    security_groups  = [aws_security_group.security_group.id]
    assign_public_ip = true
  }
}

resource "aws_security_group" "security_group" {
  description = "${var.prefix}-security-group"
  vpc_id      = var.vpc_id
}

resource "aws_security_group_rule" "ingress_rule" {
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.security_group.id
}

resource "aws_security_group_rule" "egress_rule" {
  type              = "egress"
  protocol          = "all"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.security_group.id
}
