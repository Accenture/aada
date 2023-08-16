output "s3_bucket_arn" {
  value = aws_s3_bucket.binaries_bucket.arn
}

output "ws_api_url" {
  value = "${aws_apigatewayv2_api.wsapi.execution_arn}/${aws_apigatewayv2_stage.wsapi_stage.name}/POST/@connections/*"
}

output "ws_domain_name" {
  value = aws_apigatewayv2_domain_name.wsdomain.domain_name_configuration[0].target_domain_name
}

output "http_function_url" {
  value = aws_lambda_function_url.http.function_url
}
