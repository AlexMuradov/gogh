terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }    
}

provider "aws" {
  region = "eu-central-1"
}

module "vpc" {
  source = "./modules/vpc"
  projectname = var.projectname
}

module "iam" {
  source = "./modules/iam"
  projectname = var.projectname
}

module "ecr" {
  source = "./modules/ecr"
  projectname = var.projectname
}

module "ecs" {
  source = "./modules/ecs"
  subnet_id =  module.vpc.subnet_main
  subnet_app_id =  module.vpc.subnet_app
  vpc_id = module.vpc.vpc_main
  vpc_app_id = module.vpc.vpc_app
  projectname = var.projectname
  rep_url = module.ecr.rep_url
  image = var.image
}

#module "rds" {
#  source = "./modules/rds"
#  subnet_id =  module.vpc.subnet_db
#  projectname = var.projectname
#}

# didn't have time for ci/cd, but would have done it codecommit/codebuild/codepipeline.
#module "devtools" {
#  source = "./modules/devtools"
#  vpc_id = module.vpc.vpc_main
#  subnet_id = module.vpc.subnet_main
#  security_group_id = module.vpc.sg_main
#}
