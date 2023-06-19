resource "aws_apigatewayv2_api" "httpapi" {
  name          = "${var.solution_name}-http-api"
  protocol_type = "HTTP"
  tags = {
    Name = "${var.solution_name}-http-api"
  }
}

resource "aws_apigatewayv2_route" "httpapi_default_route" {
  api_id           = aws_apigatewayv2_api.httpapi.id
  route_key        = "ANY /{proxy+}"
  api_key_required = false
  target           = "integrations/${aws_apigatewayv2_integration.httpapi_lambda.id}"
}

resource "aws_apigatewayv2_stage" "httpapi_prod_stage" {
  name        = "authenticator"
  api_id      = aws_apigatewayv2_api.httpapi.id
  auto_deploy = true

  default_route_settings {
    throttling_rate_limit  = 1000
    throttling_burst_limit = 1000
  }
}

resource "aws_apigatewayv2_integration" "httpapi_lambda" {
  api_id             = aws_apigatewayv2_api.httpapi.id
  integration_type   = "AWS_PROXY"
  integration_method = "ANY"
  integration_uri    = aws_lambda_function.http.arn
}

resource "aws_apigatewayv2_deployment" "httpapi" {
  api_id = aws_apigatewayv2_api.httpapi.id

  triggers = {
    redeployment = filesha256("../http_lambda/http_lambda.zip")
  }

  depends_on = [aws_lambda_function.http]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_domain_name" "httpdomain" {
  domain_name = "aabg.io"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.http_cert.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

output "http_endpoint" {
  value = aws_apigatewayv2_stage.httpapi_prod_stage.invoke_url
}
