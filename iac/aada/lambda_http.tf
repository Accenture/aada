resource "aws_s3_object" "http" {
  bucket      = aws_s3_bucket.code_bucket.bucket
  key         = "binaries/http_lambda.zip"
  source      = "../http_lambda/http_lambda.zip"
  source_hash = filemd5("../http_lambda/http_lambda.zip")
}

resource "aws_lambda_function" "http" {
  function_name = "${var.solution_name}-http"
  role          = var.lambda_execution_role_arn
  runtime       = "provided.al2"
  architectures = ["arm64"]
  handler       = "bootstrap"
  memory_size   = 256
  timeout       = 20 // If it doesn't happen in 20 seconds, it's not going to happen
  s3_bucket     = aws_s3_bucket.code_bucket.bucket
  s3_key        = aws_s3_object.http.key

  environment {
    variables = {
      CLIENT_ID       = var.client_id
      CLIENT_SECRET   = var.client_secret
      TABLE_NAME      = aws_dynamodb_table.data.name
      WS_CONN_URL     = "https://${aws_apigatewayv2_api.wsapi.id}.execute-api.${data.aws_region.current.name}.amazonaws.com/${aws_apigatewayv2_stage.wsapi_stage.name}/@connections"
      BINARIES_BUCKET = aws_s3_bucket.binaries_bucket.bucket
      KMS_KEY_ARN     = var.kms_key_arn
    }
  }
}

resource "aws_lambda_permission" "invoke_http_apigw" {
  statement_id_prefix = "${var.solution_name}-"
  action              = "lambda:InvokeFunction"
  function_name       = aws_lambda_function.http.function_name
  principal           = "apigateway.amazonaws.com"
  source_arn          = "${aws_apigatewayv2_api.httpapi.execution_arn}/*/*/{proxy+}"
}
//
//resource "aws_lambda_permission" "invoke_http_lex" {
//  statement_id_prefix = "${local.solution_name}-"
//  action              = "lambda:InvokeFunction"
//  function_name       = aws_lambda_function.lambda_handler.function_name
//  principal           = "lex.amazonaws.com"
//  source_arn          = "${aws_lex_intent.check_balance.id}/*"
//}