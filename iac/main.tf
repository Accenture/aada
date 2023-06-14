resource "random_pet" "solution_name" {
  // Empty
}

locals {
  solution_name       = "aada-${random_pet.solution_name.id}"
  camel_solution_name = replace(title("AADA-${random_pet.solution_name.id}"), "-", "")
}

variable "client_id" {
  type      = string
  sensitive = false
}

variable "client_secret" {
  type      = string
  sensitive = true
}

output "solution_name" {
  value = local.solution_name
}

output "camel_solution_name" {
  value = local.camel_solution_name
}

provider "aws" {
  region = "us-east-1"
  alias  = "aws_us_east_1"

  default_tags {
    tags = {
      Author      = "Eric Hill"
      Automation  = "Terraform"
      Group       = "AABG"
      Purpose     = "AADA"
      Environment = "Production"
    }
  }
}

module "aada_us_east_1" {
  source                    = "./aada"
  solution_name             = local.solution_name
  camel_solution_name       = local.camel_solution_name
  client_id                 = var.client_id
  client_secret             = var.client_secret
  lambda_execution_role_arn = aws_iam_role.lambda_execution_role.arn

  providers = {
    aws = aws.aws_us_east_1
  }
}

provider "aws" {
  region = "us-west-1"
  alias  = "aws_us_west_1"
}

module "aada_us_west_1" {
  source                    = "./aada"
  solution_name             = local.solution_name
  camel_solution_name       = local.camel_solution_name
  client_id                 = var.client_id
  client_secret             = var.client_secret
  lambda_execution_role_arn = aws_iam_role.lambda_execution_role.arn

  providers = {
    aws = aws.aws_us_west_1
  }
}