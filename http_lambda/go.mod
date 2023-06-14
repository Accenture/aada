module github.com/AccentureAWS/chitchat/http_lambda

go 1.16

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go-v2 v1.18.1
	github.com/aws/aws-sdk-go-v2/config v1.18.26
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.27
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.19.9
	github.com/aws/aws-sdk-go-v2/service/s3 v1.34.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.19.1
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/juju/ratelimit v1.0.2
	github.com/pkg/errors v0.9.1
	github.com/urfave/cli/v2 v2.2.0 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
)
