// If you want API Gateway to be able to log into CloudWatch, you need a role that API Gateway
// can assume to put logs into your account.  This really should be done at the account-level
// rather than the application level, but it's here if you want to use it.

//data "aws_iam_policy_document" "cloudwatch_assumption_policy" {
//  statement {
//    effect = "Allow"
//    principals {
//      type        = "Service"
//      identifiers = ["cloudwatch.amazonaws.com"]
//    }
//    actions = ["sts:AssumeRole"]
//  }
//}
//
//resource "aws_iam_role" "cloudwatch_logging_role" {
//  name               = "${local.solution_name}-cloudwatch-logging-role"
//  tags               = local.tags
//  assume_role_policy = data.aws_iam_policy_document.cloudwatch_assumption_policy.json
//}
//
//resource "aws_iam_role_policy_attachment" "cloudwatch_logging_role_policy" {
//  role = aws_iam_role.cloudwatch_logging_role.arn
//  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
//}
//
resource "aws_api_gateway_account" "logging" {
  cloudwatch_role_arn = "arn:aws:iam::464079168809:role/api-gateway-cloudwatch-logging"
}
