provider "aws" {
  region = "us-east-1"

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