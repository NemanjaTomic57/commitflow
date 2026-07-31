terraform {
  backend "s3" {
    bucket = "terraform-761018874759"
    key    = "commitflow-ecs.tfstate"
    region = "eu-central-1"
  }
}

##################################################
# Terraform Output Values
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
  private_subnet_ids = data.terraform_remote_state.base.outputs.private_subnet_ids
}

##################################################
# ECS Task Definitions
##################################################

module "task_definitions" {
  source = "./modules/task_definitions"

  name       = var.name
  aws_region = var.aws_region

  ecs_task_role_arn             = "arn:aws:iam::761018874759:role/ECSCommitFlowTaskRole"
  ecs_execution_role_arn        = "arn:aws:iam::761018874759:role/ECSCommitFlowTaskExecutionRole"
  ecr_commitflow_repository_url = "761018874759.dkr.ecr.eu-central-1.amazonaws.com/commitflow"
  ssm_parameter_github_pat      = "arn:aws:ssm:eu-central-1:761018874759:parameter/commitflow/passwords/github-pat"
  ssm_parameter_gitlab_pat      = "arn:aws:ssm:eu-central-1:761018874759:parameter/commitflow/passwords/gitlab-pat"
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

