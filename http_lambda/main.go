package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var kmsKeyArn string
var websocketUrl string

func main() {
	s, ok := os.LookupEnv("KMS_KEY_ARN")
	if !ok {
		fmt.Println("KMS_KEY_ARN was not provided")
	}
	kmsKeyArn = s

	s, ok = os.LookupEnv("WS_CONN_URL")
	if !ok {
		fmt.Println("WS_CONN_URL was not provided")
	}
	websocketUrl = s

	lambda.Start(lambdaHandler)
}
