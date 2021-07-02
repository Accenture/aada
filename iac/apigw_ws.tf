resource "aws_apigatewayv2_api" "wsapi" {
  name                       = "${local.solution_name}-ws-api"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"

  depends_on = [aws_api_gateway_account.logging]
}

resource "aws_apigatewayv2_route" "wsapi_default_route" {
  api_id           = aws_apigatewayv2_api.wsapi.id
  route_key        = "$default"
  api_key_required = false
  target           = "integrations/${aws_apigatewayv2_integration.wsapi_lambda.id}"
}

resource "aws_apigatewayv2_route" "wsapi_connect_route" {
  api_id           = aws_apigatewayv2_api.wsapi.id
  route_key        = "$connect"
  api_key_required = false
  target           = "integrations/${aws_apigatewayv2_integration.wsapi_lambda.id}"
}

resource "aws_apigatewayv2_route" "wsapi_disconnect_route" {
  api_id           = aws_apigatewayv2_api.wsapi.id
  route_key        = "$disconnect"
  api_key_required = false
  target           = "integrations/${aws_apigatewayv2_integration.wsapi_lambda.id}"
}

resource "aws_apigatewayv2_route_response" "wsapi_default_route" {
  api_id             = aws_apigatewayv2_api.wsapi.id
  route_id           = aws_apigatewayv2_route.wsapi_default_route.id
  route_response_key = "$default"
}

resource "aws_apigatewayv2_stage" "wsapi_stage" {
  name        = "chat"
  api_id      = aws_apigatewayv2_api.wsapi.id
  auto_deploy = true

  default_route_settings {
    throttling_rate_limit    = 1000
    throttling_burst_limit   = 1000
    data_trace_enabled       = true
    detailed_metrics_enabled = true
    logging_level            = "INFO"
  }
}

resource "aws_apigatewayv2_integration" "wsapi_lambda" {
  api_id                    = aws_apigatewayv2_api.wsapi.id
  integration_type          = "AWS_PROXY"
  connection_type           = "INTERNET"
  content_handling_strategy = "CONVERT_TO_TEXT"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.ws.invoke_arn
  passthrough_behavior      = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_integration_response" "wsapi_lambda" {
  api_id                   = aws_apigatewayv2_api.wsapi.id
  integration_id           = aws_apigatewayv2_integration.wsapi_lambda.id
  integration_response_key = "$default"
}

resource "aws_apigatewayv2_deployment" "wsapi" {
  api_id = aws_apigatewayv2_api.wsapi.id

  triggers = {
    redeployment = filesha256("../ws_lambda/ws_lambda.zip")
  }

  lifecycle {
    create_before_destroy = true
  }
}

output "ws_endpoint" {
  value = aws_apigatewayv2_stage.wsapi_stage.invoke_url
}