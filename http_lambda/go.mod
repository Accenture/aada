module github.com/AccentureAWS/chitchat/http_lambda

go 1.22

toolchain go1.23.1

require (
	github.com/aws/aws-lambda-go v1.49.0
	github.com/aws/aws-sdk-go-v2 v1.38.3
	github.com/aws/aws-sdk-go-v2/config v1.31.6
	github.com/aws/aws-sdk-go-v2/service/kms v1.45.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.87.3
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.2
	github.com/fxamacker/cbor/v2 v2.9.0
	github.com/google/uuid v1.6.0
	github.com/juju/ratelimit v1.0.2
	github.com/pkg/errors v0.9.1
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.8.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.34.2 // indirect
	github.com/aws/smithy-go v1.23.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 // indirect
)
