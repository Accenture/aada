resource "aws_s3_bucket" "code_bucket" {
  bucket = "${local.solution_name}-code-bucket"
}
