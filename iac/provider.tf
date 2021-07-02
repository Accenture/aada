provider "aws" {
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