resource "aws_s3_bucket" "code_bucket" {
  bucket = "${local.solution_name}-code-bucket"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}
