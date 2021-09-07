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

resource "aws_s3_bucket" "binaries_bucket" {
  bucket = "${local.solution_name}-binaries"
  acl    = "public-read"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}

resource "aws_s3_bucket_object" "binaries" {
  bucket   = aws_s3_bucket.binaries_bucket.bucket
  for_each = fileset("../binaries", "*.zip")
  acl      = "public-read"
  key      = each.value
  source   = "../binaries/${each.value}"
}
