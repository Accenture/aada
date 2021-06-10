resource "random_pet" "solution_name" {
  // Empty
}

locals {
  solution_name       = "aada-${random_pet.solution_name.id}"
  camel_solution_name = replace(title("AADA-${random_pet.solution_name.id}"), "-", "")
  tags = {
    Author = "Eric Hill"
  }
}

data "aws_region" "current" {
  // Empty
}

output "solution_name" {
  value = local.solution_name
}

output "camel_solution_name" {
  value = local.camel_solution_name
}
