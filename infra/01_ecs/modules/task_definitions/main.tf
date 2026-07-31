##################################################
# VPC Data
##################################################

data "terraform_remote_state" "base" {
  backend = "s3"

  config = {
    bucket = "terraform-761018874759"
    key    = "commitflow-base.tfstate"
    region = "eu-central-1"
  }
}

locals {
  ssm_parameter_db_engine             = data.terraform_remote_state.base.outputs.ssm_parameter_db_engine
  ssm_parameter_db_address            = data.terraform_remote_state.base.outputs.ssm_parameter_db_address
  ssm_parameter_db_port               = data.terraform_remote_state.base.outputs.ssm_parameter_db_port
  ssm_parameter_db_name               = data.terraform_remote_state.base.outputs.ssm_parameter_db_name
  ssm_parameter_db_username           = data.terraform_remote_state.base.outputs.ssm_parameter_db_username
  db_secret_arn                       = data.terraform_remote_state.base.outputs.db_secret_arn
  ssm_parameter_kafka_bootstrap_sever = data.terraform_remote_state.base.outputs.ssm_parameter_kafka_bootstrap_server
}

##################################################
# Cloud Watch
##################################################

resource "aws_cloudwatch_log_group" "this" {
  name = "/ecs/${var.name}"
}

##################################################
# ECS Task Definitions
##################################################

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
          valueFrom = local.ssm_parameter_kafka_bootstrap_sever
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

resource "aws_ecs_task_definition" "commitflow_consumer" {
  family                   = "commitflow-consumer"
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
      name      = "consumer"
      image     = "${var.ecr_commitflow_repository_url}:latest"
      essential = true

      command = ["consumer"]

      secrets = [
        {
          name      = "DB_ENGINE"
          valueFrom = local.ssm_parameter_db_engine
        },
        {
          name      = "DB_ADDRESS"
          valueFrom = local.ssm_parameter_db_address
        },
        {
          name      = "DB_PORT"
          valueFrom = local.ssm_parameter_db_port
        },
        {
          name      = "DB_NAME"
          valueFrom = local.ssm_parameter_db_name
        },
        {
          name      = "DB_USERNAME"
          valueFrom = local.ssm_parameter_db_username
        },
        {
          name      = "DB_PASSWORD"
          valueFrom = local.db_secret_arn
        },
        {
          name      = "KAFKA_BOOTSTRAP_SERVER"
          valueFrom = local.ssm_parameter_kafka_bootstrap_sever
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
          awslogs-stream-prefix = "consumer"
        }
      }
    }
  ])

  tags = {
    Name = "${var.name}-commitflow-consumer"
  }
}
