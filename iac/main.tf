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

resource "aws_kms_key" "signatory" {
  provider                 = aws.aws_us_east_1
  description              = "AADA Signatory"
  key_usage                = "SIGN_VERIFY"
  customer_master_key_spec = "ECC_NIST_P256"
  multi_region             = true
  deletion_window_in_days  = 10
}


module "aada_us_east_1" {
  source                    = "./aada"
  solution_name             = local.solution_name
  camel_solution_name       = local.camel_solution_name
  client_id                 = var.client_id
  client_secret             = var.client_secret
  lambda_execution_role_arn = aws_iam_role.lambda_execution_role.arn
  kms_key_arn               = aws_kms_key.signatory.arn

  providers = {
    aws = aws.aws_us_east_1
  }
}

provider "aws" {
  region = "us-west-1"
  alias  = "aws_us_west_1"

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

moved {
  from = aws_kms_replica_key.signatory
  to   = aws_kms_replica_key.signatory_us_west_1
}

resource "aws_kms_replica_key" "signatory_us_west_1" {
  provider                = aws.aws_us_west_1
  primary_key_arn         = aws_kms_key.signatory.arn
  description             = aws_kms_key.signatory.description
  deletion_window_in_days = aws_kms_key.signatory.deletion_window_in_days
}

module "aada_us_west_1" {
  source                    = "./aada"
  solution_name             = local.solution_name
  camel_solution_name       = local.camel_solution_name
  client_id                 = var.client_id
  client_secret             = var.client_secret
  lambda_execution_role_arn = aws_iam_role.lambda_execution_role.arn
  kms_key_arn               = aws_kms_replica_key.signatory_us_west_1.arn

  providers = {
    aws = aws.aws_us_west_1
  }
}

provider "aws" {
  region = "eu-central-1"
  alias  = "aws_eu_central_1"

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

resource "aws_kms_replica_key" "signatory_eu_central_1" {
  provider                = aws.aws_eu_central_1
  primary_key_arn         = aws_kms_key.signatory.arn
  description             = aws_kms_key.signatory.description
  deletion_window_in_days = aws_kms_key.signatory.deletion_window_in_days
}

module "aada_eu_central_1" {
  source                    = "./aada"
  solution_name             = local.solution_name
  camel_solution_name       = local.camel_solution_name
  client_id                 = var.client_id
  client_secret             = var.client_secret
  lambda_execution_role_arn = aws_iam_role.lambda_execution_role.arn
  kms_key_arn               = aws_kms_replica_key.signatory_eu_central_1.arn

  providers = {
    aws = aws.aws_eu_central_1
  }
}

provider "aws" {
  region = "ap-southeast-2"
  alias  = "aws_ap_southeast_2"

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

resource "aws_kms_replica_key" "signatory_ap_southeast_2" {
  provider                = aws.aws_ap_southeast_2
  primary_key_arn         = aws_kms_key.signatory.arn
  description             = aws_kms_key.signatory.description
  deletion_window_in_days = aws_kms_key.signatory.deletion_window_in_days
}

module "aada_ap_southeast_2" {
  source                    = "./aada"
  solution_name             = local.solution_name
  camel_solution_name       = local.camel_solution_name
  client_id                 = var.client_id
  client_secret             = var.client_secret
  lambda_execution_role_arn = aws_iam_role.lambda_execution_role.arn
  kms_key_arn               = aws_kms_replica_key.signatory_ap_southeast_2.arn

  providers = {
    aws = aws.aws_ap_southeast_2
  }
}