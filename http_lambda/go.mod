module github.com/AccentureAWS/chitchat/http_lambda

go 1.16

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go-v2 v1.22.2
	github.com/aws/aws-sdk-go-v2/config v1.24.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.30
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.20.0
	github.com/aws/aws-sdk-go-v2/service/kms v1.26.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.42.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.25.1
	github.com/fxamacker/cbor/v2 v2.5.0
	github.com/google/uuid v1.4.0
	github.com/juju/ratelimit v1.0.2
	github.com/pkg/errors v0.9.1
)
