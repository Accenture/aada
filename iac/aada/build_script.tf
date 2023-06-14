#resource "local_file" "deploy_ws" {
#  filename = "../deploy-ws.sh"
#  content  = <<-EOF
#  #!/bin/bash
#
#  aws s3 cp ws_lambda/ws_lambda.zip s3://${aws_s3_bucket.code_bucket.bucket}/${aws_s3_object.ws.key}
#  aws lambda update-function-code --function-name ${aws_lambda_function.ws.function_name} --s3-bucket ${aws_s3_bucket.code_bucket.bucket} --s3-key ${aws_s3_object.ws.key}
#  EOF
#}
#
#resource "local_file" "deploy_http" {
#  filename = "../deploy-http.sh"
#  content  = <<-EOF
#  #!/bin/bash
#
#  aws s3 cp http_lambda/http_lambda.zip s3://${aws_s3_bucket.code_bucket.bucket}/${aws_s3_object.http.key}
#  aws lambda update-function-code --function-name ${aws_lambda_function.http.function_name} --s3-bucket ${aws_s3_bucket.code_bucket.bucket} --s3-key ${aws_s3_object.http.key}
#  EOF
#}
