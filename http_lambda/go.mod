module github.com/AccentureAWS/chitchat/http_lambda

go 1.16

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go-v2 v1.18.1
	github.com/aws/aws-sdk-go-v2/config v1.18.27
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.30
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.20.0
	github.com/aws/aws-sdk-go-v2/service/kms v1.23.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.36.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.19.2
	github.com/fxamacker/cbor/v2 v2.4.0
	github.com/google/uuid v1.3.0
	github.com/juju/ratelimit v1.0.2
	github.com/pkg/errors v0.9.1
)
