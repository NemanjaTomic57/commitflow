terraform {
  backend "s3" {
    bucket = "terraform-761018874759"
    key    = "commitflow-base.tfstate"
    region = "eu-central-1"
  }
}

locals {
  name = "commitflow"
}

module "vpc" {
  source = "./modules/vpc"

  name = local.name

  vpc_cidr = "10.100.0.0/16"

  public_subnet_cidrs = [
    "10.100.1.0/24",
    "10.100.2.0/24",
    "10.100.3.0/24"
  ]

  private_subnet_cidrs = [
    "10.100.16.0/24",
    "10.100.17.0/24",
    "10.100.18.0/24"
  ]
}

module "kafka" {
  source = "./modules/kafka"

  name               = local.name
  ami_id             = "ami-0723bff07f72bb394"
  key_name           = "aws"
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  private_subnet_ids = module.vpc.private_subnet_ids
}
