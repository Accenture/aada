resource "aws_s3_bucket" "code_bucket" {
  bucket = "${var.solution_name}-${data.aws_region.current.name}-code-bucket"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "code_bucket" {
  bucket = aws_s3_bucket.code_bucket.bucket
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "code_bucket" {
  bucket              = aws_s3_bucket.code_bucket.bucket
  block_public_acls   = true
  block_public_policy = true
}

resource "aws_s3_bucket" "binaries_bucket" {
  bucket = "${var.solution_name}-${data.aws_region.current.name}-binaries"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "binaries_bucket" {
  bucket = aws_s3_bucket.binaries_bucket.bucket
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "binaries_block" {
  bucket              = aws_s3_bucket.binaries_bucket.bucket
  block_public_acls   = true
  block_public_policy = true
}
