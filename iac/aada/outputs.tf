output "dynamodb_table_arn" {
  value = aws_dynamodb_table.data.arn
}

output "s3_bucket_arn" {
  value = aws_s3_bucket.binaries_bucket.arn
}

output "ws_api_url" {
  value = "${aws_apigatewayv2_api.wsapi.execution_arn}/${aws_apigatewayv2_stage.wsapi_stage.name}/POST/@connections/*"
}