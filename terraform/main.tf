terraform {
  required_version = ">= 1.0.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }

  backend "s3" {
    bucket         = "gitsetrun-terraform-remote-state"
    key            = "terraform/gitsetrun/terraform.tfstate"
    region         = "us-west-2"
    dynamodb_table = "gitsetrun-terraform-state-locks"
    encrypt        = true
  }
}

provider "aws" {
  region = "us-west-2"
}
