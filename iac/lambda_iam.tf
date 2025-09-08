data "aws_iam_policy_document" "lambda_assumption_policy" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda_execution_role" {
  provider           = aws.aws_us_east_1
  name               = "aada-trustpoint"
  assume_role_policy = data.aws_iam_policy_document.lambda_assumption_policy.json

  lifecycle {
    prevent_destroy = true // If this role gets destroyed, we have a big problem
  }
}

data "aws_iam_policy_document" "lambda_policy" {
  statement {
    sid    = "KmsSignVerify"
    effect = "Allow"
    actions = [
      "kms:Sign",
      "kms:Verify"
    ]
    resources = ["*"]
  }
  statement {
    sid       = "DownloadClientBinary"
    effect    = "Allow"
    actions   = ["s3:GetObject"]
    resources = ["${module.aada_us_east_1.s3_bucket_arn}/*", "${module.aada_us_west_1.s3_bucket_arn}/*"]
  }
  statement {
    sid     = "WSSAsyncPush"
    effect  = "Allow"
    actions = ["execute-api:ManageConnections"] // https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-websocket-control-access-iam.html
    resources = [
      module.aada_us_east_1.ws_api_url,
      module.aada_us_west_1.ws_api_url,
      module.aada_ap_southeast_2.ws_api_url,
      module.aada_eu_central_1.ws_api_url
    ]
  }
  statement {
    sid       = "CrossAccountAssumptions"
    effect    = "Allow"
    actions   = ["sts:AssumeRole"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "lambda_role_policy" {
  provider    = aws.aws_us_east_1
  name_prefix = local.solution_name
  role        = aws_iam_role.lambda_execution_role.id
  policy      = data.aws_iam_policy_document.lambda_policy.json
}

resource "aws_iam_role_policy_attachment" "lambda_role_basics" {
  provider   = aws.aws_us_east_1
  role       = aws_iam_role.lambda_execution_role.id
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
