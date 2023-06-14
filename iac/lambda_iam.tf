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
  name               = "aada-trustpoint"
  assume_role_policy = data.aws_iam_policy_document.lambda_assumption_policy.json
}

data "aws_iam_policy_document" "lambda_policy" {
  statement {
    sid    = "DynamoDataAccess"
    effect = "Allow"
    actions = [
      "dynamodb:Batch*", // Batch record manipulation
      "dynamodb:ConditionCheckItem",
      "dynamodb:DeleteItem",
      "dynamodb:Describe*", // Descriptive access
      "dynamodb:Get*",      // Read-only access
      "dynamodb:PartiQL*",  // PartiQL full access
      "dynamodb:PutItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:UpdateItem"
    ]
    resources = [module.aada_us_east_1.dynamodb_table_arn]
  }
  statement {
    sid       = "DownloadClientBinary"
    effect    = "Allow"
    actions   = ["s3:GetObject"]
    resources = ["${module.aada_us_east_1.s3_bucket_arn}/*"]
  }
  statement {
    sid       = "WSSAsyncPush"
    effect    = "Allow"
    actions   = ["execute-api:ManageConnections"] // https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-websocket-control-access-iam.html
    resources = [module.aada_us_east_1.ws_api_url]
  }
  statement {
    sid       = "CrossAccountAssumptions"
    effect    = "Allow"
    actions   = ["sts:AssumeRole"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "lambda_role_policy" {
  name_prefix = local.solution_name
  role        = aws_iam_role.lambda_execution_role.id
  policy      = data.aws_iam_policy_document.lambda_policy.json
}

resource "aws_iam_role_policy_attachment" "lambda_role_basics" {
  role       = aws_iam_role.lambda_execution_role.id
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
