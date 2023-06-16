package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var kmsKeyArn string

func main() {
	s, ok := os.LookupEnv("KMS_KEY_ARN")
	if !ok {
		fmt.Println("KMS_KEY_ARN was not provided")
	}
	kmsKeyArn = s

	lambda.Start(lambdaHandler)
}
