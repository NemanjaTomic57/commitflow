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

  public_subnet_cidrs = {
    "public-1" = {
      cidr = "10.100.1.0/24",
      az   = "eu-central-1a"
    },
    "public-2" = {
      cidr = "10.100.2.0/24",
      az   = "eu-central-1b"
    },
    "public-3" = {
      cidr = "10.100.3.0/24",
      az   = "eu-central-1c"
    },
  }

  private_subnet_cidrs = {
    "private-1" = {
      cidr = "10.100.16.0/24",
      az   = "eu-central-1a"
    },
    "private-2" = {
      cidr = "10.100.17.0/24",
      az   = "eu-central-1b"
    },
    "private-3" = {
      cidr = "10.100.18.0/24",
      az   = "eu-central-1c"
    },
  }
}

module "security" {
  source = "./modules/security"

  name     = local.name
  vpc_id   = module.vpc.vpc_id
  vpc_cidr = module.vpc.vpc_cidr
}

module "rds" {
  source = "./modules/rds"

  name                 = local.name
  instance_class       = "db.t4g.micro"
  storage_type         = "gp3"
  allocated_storage    = 20
  private_subnet_ids   = module.vpc.private_subnet_ids
  db_security_group_id = module.security.db_security_group_id
}

module "nat" {
  source = "./modules/nat"

  name                  = local.name
  vpc_id                = module.vpc.vpc_id
  internet_gateway_id   = module.vpc.internet_gateway_id
  public_subnet_ids     = module.vpc.public_subnet_ids
  private_subnet_ids    = module.vpc.private_subnet_ids
  nat_security_group_id = module.security.nat_security_group_id
  ami_id                = "ami-0d48b1c648cd339e0"
  key_name              = "aws"
}

module "kafka" {
  source = "./modules/kafka"

  name                    = local.name
  vpc_id                  = module.vpc.vpc_id
  private_subnet_ids      = module.vpc.private_subnet_ids
  kafka_security_group_id = module.security.kafka_security_group_id
  ami_id                  = "ami-0723bff07f72bb394"
  key_name                = "aws"
}
