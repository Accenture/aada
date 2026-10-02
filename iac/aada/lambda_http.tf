resource "aws_s3_object" "http" {
  bucket      = aws_s3_bucket.code_bucket.bucket
  key         = "binaries/http_lambda.zip"
  source      = "../http_lambda/http_lambda.zip"
  source_hash = filemd5("../http_lambda/http_lambda.zip")
}

resource "aws_secretsmanager_secret" "client_secret" {
  name                    = "${var.solution_name}-client-secret"
  description             = "Azure AD client secret for AADA"
  recovery_window_in_days = 7
}

resource "aws_secretsmanager_secret_version" "client_secret" {
  secret_id     = aws_secretsmanager_secret.client_secret.id
  secret_string = var.client_secret
}

resource "aws_lambda_function" "http" {
  function_name    = "${var.solution_name}-http"
  role             = var.lambda_execution_role_arn
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  handler          = "bootstrap"
  memory_size      = 256
  timeout          = 20 // If it doesn't happen in 20 seconds, it's not going to happen
  s3_bucket        = aws_s3_bucket.code_bucket.bucket
  s3_key           = aws_s3_object.http.key
  source_code_hash = filebase64sha256(aws_s3_object.http.source)

  environment {
    variables = {
      CLIENT_ID          = var.client_id
      CLIENT_SECRET_ARN  = aws_secretsmanager_secret.client_secret.arn
      WS_CONN_URL        = "https://${aws_apigatewayv2_api.wsapi.id}.execute-api.${data.aws_region.current.region}.amazonaws.com/${aws_apigatewayv2_stage.wsapi_stage.name}/@connections"
      BINARIES_BUCKET    = aws_s3_bucket.binaries_bucket.bucket
      KMS_KEY_ARN        = var.kms_key_arn
    }
  }
}

resource "aws_lambda_function_url" "http" {
  function_name      = aws_lambda_function.http.function_name
  authorization_type = "NONE"
}
