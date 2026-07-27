terraform {
  backend "s3" {
    bucket = "terraform-761018874759"
    key    = "commitflow-ecs.tfstate"
    region = "eu-central-1"
  }
}

data "terraform_remote_state" "vpc" {
  backend = "s3"

  config = {
    bucket = "terraform-761018874759"
    key    = "commitflow-base.tfstate"
    region = "eu-central-1"
  }
}

resource "aws_ecs_cluster" "commitflow" {
  name   = var.name
  region = var.aws_region
}

module "task_definitions" {
  source = "./modules/task_definitions"

  name                          = var.name
  aws_region                    = var.aws_region
  kafka_bootstrap_server        = data.terraform_remote_state.vpc.outputs.kafka_bootstrap_server
  ecr_commitflow_repository_url = "761018874759.dkr.ecr.eu-central-1.amazonaws.com/commitflow"
}
