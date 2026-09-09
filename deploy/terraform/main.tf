terraform {
  required_version = ">= 1.0"
}

provider "aws" {
  region = var.aws_region
}

module "vpc" {
  source     = "./modules/vpc"
  cidr_block = var.vpc_cidr
}

module "lab" {
  source    = "./modules/lab"
  vpc_id    = module.vpc.vpc_id
  subnet_id = module.vpc.subnet_id
}
