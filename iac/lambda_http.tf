resource "aws_s3_object" "http" {
  bucket = aws_s3_bucket.code_bucket.bucket
  key    = "binaries/http_lambda.zip"
  source = "../http_lambda/http_lambda.zip"
}

variable "client_secret" {
  type      = string
  sensitive = true
}

resource "aws_lambda_function" "http" {
  function_name = "${local.solution_name}-http"
  role          = aws_iam_role.lambda_execution_role.arn
  runtime       = "go1.x"
  handler       = "http_lambda"
  memory_size   = 256
  timeout       = 60
  s3_bucket     = aws_s3_bucket.code_bucket.bucket
  s3_key        = aws_s3_object.http.key

  environment {
    variables = {
      CLIENT_SECRET = var.client_secret
      TABLE_NAME    = aws_dynamodb_table.data.name
      WS_CONN_URL   = "https://${aws_apigatewayv2_api.wsapi.id}.execute-api.${data.aws_region.current.name}.amazonaws.com/${aws_apigatewayv2_stage.wsapi_stage.name}/@connections"
    }
  }
}

resource "aws_lambda_permission" "invoke_http_apigw" {
  statement_id_prefix = "${local.solution_name}-"
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
