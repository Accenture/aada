module github.com/AccentureAWS/chitchat/http_lambda

go 1.16

require (
	github.com/aws/aws-lambda-go v1.46.0
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.26.6
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.30
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.20.0
	github.com/aws/aws-sdk-go-v2/service/kms v1.27.9
	github.com/aws/aws-sdk-go-v2/service/s3 v1.48.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.7
	github.com/fxamacker/cbor/v2 v2.5.0
	github.com/google/uuid v1.6.0
	github.com/juju/ratelimit v1.0.2
	github.com/pkg/errors v0.9.1
)
