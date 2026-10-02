package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var kmsKeyArn string
var clientSecret string

func main() {
	s, ok := os.LookupEnv("KMS_KEY_ARN")
	if !ok {
		fmt.Println("KMS_KEY_ARN was not provided")
	}
	kmsKeyArn = s

	secretArn, ok := os.LookupEnv("CLIENT_SECRET_ARN")
	if !ok {
		fmt.Println("ERROR environment variable CLIENT_SECRET_ARN was not provided")
	} else {
		smc := secretsmanager.NewFromConfig(awsConfig)
		out, err := smc.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
			SecretId: &secretArn,
		})
		if err != nil {
			fmt.Println("ERROR fetching client secret from Secrets Manager: ", err.Error())
		} else if out.SecretString != nil {
			clientSecret = *out.SecretString
		}
	}

	lambda.Start(lambdaHandler)
}
