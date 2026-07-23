terraform {
  backend "s3" {
    bucket = "terraform-761018874759"
    key    = "commitflow.tfstate"
    region = "eu-central-1"
  }
}

module "vpc" {
  source = "./modules/vpc"

  name = "commitflow"

  vpc_cidr = "10.100.0.0/16"

  public_subnet_cidrs = [
    "10.100.1.0/24",
    "10.100.2.0/24"
  ]

  private_subnet_cidrs = [
    "10.100.16.0/24",
    "10.100.17.0/24"
  ]
}

module "kafka_cluster" {
  source =
}
