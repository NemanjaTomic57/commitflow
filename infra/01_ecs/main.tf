terraform {
  backend "s3" {
    bucket = "terraform-761018874759"
    key    = "commitflow-ecs.tfstate"
    region = "eu-central-1"
  }
}

##################################################
# VPC Data
##################################################

data "terraform_remote_state" "vpc" {
  backend = "s3"

  config = {
    bucket = "terraform-761018874759"
    key    = "commitflow-base.tfstate"
    region = "eu-central-1"
  }
}

locals {
  kafka_bootstrap_server = data.terraform_remote_state.vpc.outputs.kafka_bootstrap_server
  private_subnet_ids     = data.terraform_remote_state.vpc.outputs.private_subnet_ids
}

##################################################
# ECS Task Definitions
##################################################

module "task_definitions" {
  source = "./modules/task_definitions"

  name                          = var.name
  aws_region                    = var.aws_region
  kafka_bootstrap_server        = local.kafka_bootstrap_server
  ecr_commitflow_repository_url = "761018874759.dkr.ecr.eu-central-1.amazonaws.com/commitflow"
}

##################################################
# ECS Cluster
##################################################

resource "aws_ecs_cluster" "commitflow" {
  name   = var.name
  region = var.aws_region
}

##################################################
# ECS Services
##################################################

resource "aws_ecs_service" "producer" {
  name                 = "producer"
  cluster              = aws_ecs_cluster.commitflow.id
  task_definition      = module.task_definitions.producer_task_definition_arn
  desired_count        = 1
  launch_type          = "FARGATE"
  force_new_deployment = true

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  network_configuration {
    subnets = values(local.private_subnet_ids)
  }
}

