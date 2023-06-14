resource "aws_s3_object" "ws" {
  bucket = aws_s3_bucket.code_bucket.bucket
  key    = "binaries/ws_lambda.zip"
  source = "../ws_lambda/ws_lambda.zip"
}

resource "aws_lambda_function" "ws" {
  function_name = "${var.solution_name}-ws"
  role          = var.lambda_execution_role_arn
  runtime       = "go1.x"
  handler       = "ws_lambda"
  memory_size   = 256
  timeout       = 10
  s3_bucket     = aws_s3_bucket.code_bucket.bucket
  s3_key        = aws_s3_object.ws.key

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.data.name
    }
  }
}

resource "aws_lambda_permission" "invoke_ws" {
  statement_id_prefix = "${var.solution_name}-"
  action              = "lambda:InvokeFunction"
  function_name       = aws_lambda_function.ws.function_name
  principal           = "apigateway.amazonaws.com"
  source_arn          = "${aws_apigatewayv2_api.wsapi.execution_arn}/*"
}
