resource "aws_dynamodb_table" "data" {
  name         = "${var.solution_name}-data"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "state"
  ttl {
    attribute_name = "expiration"
    enabled        = true
  }
  attribute {
    name = "state"
    type = "S"
  }
}