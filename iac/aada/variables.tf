variable "client_id" {
  type      = string
  sensitive = false
}

variable "client_secret" {
  type      = string
  sensitive = true
}

variable "solution_name" {
  type = string
}

variable "camel_solution_name" {
  type = string
}

variable "lambda_execution_role_arn" {
  type = string
}