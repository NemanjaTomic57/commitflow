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
  ssm_parameter_db_password           = data.terraform_remote_state.base.outputs.ssm_parameter_db_password
  ssm_parameter_kafka_bootstrap_sever = data.terraform_remote_state.base.outputs.ssm_parameter_kafka_bootstrap_server
  ssm_parameter_github_pat            = data.terraform_remote_state.base.outputs.ssm_parameter_github_pat
  ssm_parameter_gitlab_pat            = data.terraform_remote_state.base.outputs.ssm_parameter_gitlab_pat
  ssm_parameter_grafana_username      = data.terraform_remote_state.base.outputs.ssm_parameter_grafana_username
  ssm_parameter_grafana_password      = data.terraform_remote_state.base.outputs.ssm_parameter_grafana_password
}

##################################################
# Cloud Watch
##################################################

resource "aws_cloudwatch_log_group" "this" {
  name = "/ecs/${var.name}"
}

##################################################
# ECS Task Definition CommitFlow Producer
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
          valueFrom = local.ssm_parameter_github_pat
        },
        {
          name      = "GITLAB_PAT"
          valueFrom = local.ssm_parameter_gitlab_pat
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

##################################################
# ECS Task Definition CommitFlow Consumer
##################################################

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
          name      = "DB_USER"
          valueFrom = local.ssm_parameter_db_username
        },
        {
          name      = "DB_PASSWORD"
          valueFrom = local.ssm_parameter_db_password
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

##################################################
# ECS Task Definition Grafana
##################################################

resource "aws_ecs_task_definition" "grafana" {
  family                   = "grafana"
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
      name      = "grafana"
      image     = "grafana/grafana:12.2"
      essential = true

      secrets = [
        {
          name      = "GF_SECURITY_ADMIN_USER"
          valueFrom = local.ssm_parameter_grafana_username
        },
        {
          name      = "GF_SECURITY_ADMIN_PASSWORD"
          valueFrom = local.ssm_parameter_grafana_password
        },
        {
          name      = "GF_DATABASE_TYPE"
          valueFrom = local.ssm_parameter_db_engine
        },
        {
          name      = "GF_DATABASE_HOST"
          valueFrom = local.ssm_parameter_db_address
        },
        {
          name      = "GF_DATABASE_NAME"
          valueFrom = local.ssm_parameter_db_name
        },
        {
          name      = "GF_DATABASE_USER"
          valueFrom = local.ssm_parameter_db_username
        },
        {
          name      = "GF_DATABASE_PASSWORD"
          valueFrom = local.ssm_parameter_db_password
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
          awslogs-stream-prefix = "grafana"
        }
      }
    }
  ])

  tags = {
    Name = "${var.name}-commitflow-grafana"
  }
}
