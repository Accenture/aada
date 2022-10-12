resource "aws_cloudwatch_log_group" "apigw_ws" {
  name              = "/aws/apigateway/${aws_apigatewayv2_api.wsapi.id}/${aws_apigatewayv2_stage.wsapi_stage.name}"
  retention_in_days = 3
}

resource "aws_cloudwatch_log_group" "apigw_http" {
  name              = "/aws/apigateway/${aws_apigatewayv2_api.httpapi.id}/${aws_apigatewayv2_stage.httpapi_prod_stage.name}"
  retention_in_days = 30
}

resource "aws_cloudwatch_log_group" "lambda_ws" {
  name              = "/aws/lambda/${aws_lambda_function.ws.function_name}"
  retention_in_days = 3
}

resource "aws_cloudwatch_log_group" "lambda_http" {
  name              = "/aws/lambda/${aws_lambda_function.http.function_name}"
  retention_in_days = 3
}

resource "aws_cloudwatch_log_metric_filter" "throttle_filter" {
  name           = "${local.solution_name}-throttle-filter"
  pattern        = "THROTTLE"
  log_group_name = aws_cloudwatch_log_group.lambda_http.name

  metric_transformation {
    name      = "ThrottleCount"
    namespace = local.camel_solution_name
    value     = "1"
  }
}