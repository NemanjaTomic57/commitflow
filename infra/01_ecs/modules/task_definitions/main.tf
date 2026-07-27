resource "aws_ssm_parameter" "kafka_bootstrap_server" {
  name  = "/${var.name}/kafka/bootstrap-server"
  type  = "String"
  value = var.kafka_bootstrap_server
}

resource "aws_cloudwatch_log_group" "this" {
  name = "/ecs/${var.name}"
}

resource "aws_ecs_task_definition" "commitflow_producer" {
  family                   = "commitflow-producer"
  requires_compatibilities = ["FARGATE"]

  network_mode = "awsvpc"
  cpu          = "256"
  memory       = "512"

  task_role_arn      = var.ecs_task_role_arn
  execution_role_arn = var.ecs_execution_role_arn

  runtime_platform {
    operating_system_family = "LINUX"
  }

  pid_mode = "task"

  container_definitions = jsonencode([
    {
      name      = "producer"
      image     = "${var.ecr_commitflow_repository_url}:latest"
      essential = true

      command = [
        "producer",
        "-bootstrap"
      ]

      secrets = [
        {
          name      = "KAFKA_BOOTSTRAP_SERVER"
          valueFrom = aws_ssm_parameter.kafka_bootstrap_server.arn
        },
        {
          name      = "GITHUB_PAT"
          valueFrom = var.ssm_parameter_github_pat
        },
        {
          name      = "GITLAB_PAT"
          valueFrom = var.ssm_parameter_gitlab_pat
        }
      ]

      linuxParameters = {
        initProcessEnabled = true
      }

      logConfiguration = {
        logDriver = "awslogs"

        options = {
          awslogs-group         = aws_cloudwatch_log_group.this.name
          awslogs-region        = var.aws_region
          awslogs-stream-prefix = "producer"
        }
      }
    }
  ])

  tags = {
    Name = "${var.name}-commitflow-producer"
  }
}
