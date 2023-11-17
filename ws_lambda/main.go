package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"net/http"
	"os"
)

var kmsKeyArn string
var awsRegion string

func main() {
	s, ok := os.LookupEnv("KMS_KEY_ARN")
	if !ok {
		fmt.Println("KMS_KEY_ARN was not provided")
	}
	kmsKeyArn = s

	// https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html
	s, ok = os.LookupEnv("AWS_DEFAULT_REGION")
	if !ok {
		fmt.Println("AWS_DEFAULT_REGION was not provided")
	}
	awsRegion = s

	// If defined, this value overrides the AWS_DEFAULT_REGION.
	s, ok = os.LookupEnv("AWS_REGION")
	if ok {
		awsRegion = s
	}

	lambda.Start(lambdaHandler)
}

type Frame struct {
	Nonce         string `json:"nonce"`
	Profile       string `json:"profile"`
	State         string `json:"state"`
	Mode          string `json:"mode"`
	ClientVersion string `json:"client_version"`
	Duration      int    `json:"duration"`
	Connection    string `json:"-"`
}

func processMessage(ctx context.Context, event Event) HTTPResponse {
	frame := &Frame{
		Connection: event.Context.ConnectionId,
	}
	err := json.Unmarshal([]byte(event.Body), frame)
	if err != nil {
		fmt.Println("error unmarshalling body")
		return HTTPResponse{StatusCode: http.StatusInternalServerError}
	}
	if len(frame.Nonce) == 0 || len(frame.Profile) == 0 {
		fmt.Println("error with initial request")
		return HTTPResponse{StatusCode: http.StatusBadRequest}
	}

	frame.State = uuid.NewString()

	// Package up the information the caller needs to carry into a signed structure
	info := &Information{
		ConnectionId:  event.Context.ConnectionId,
		ApiId:         event.Context.ApiId,
		ProfileName:   frame.Profile,
		ConnectMode:   ModeUnknown,
		ClientVersion: frame.ClientVersion,
		AWSRegion:     awsRegion,
		Duration:      frame.Duration,
	}
	switch frame.Mode {
	case "access":
		info.ConnectMode = ModeAccess
	case "configuration":
		info.ConnectMode = ModeConfiguration
	}
	signed, err := info.Sign(ctx)
	if err != nil {
		fmt.Println("error signing connection id", err.Error())
		// Fast return the old method
		return HTTPResponse{
			StatusCode: http.StatusOK,
			Body:       fmt.Sprintf("{\"state\":\"%s\"}", frame.State),
		}
	}
	signature, err := signed.EncodeToString()
	if err != nil {
		fmt.Println("error signing connection id", err.Error())
		// Fast return the old method
		return HTTPResponse{
			StatusCode: http.StatusOK,
			Body:       fmt.Sprintf("{\"state\":\"%s\"}", frame.State),
		}
	}
	fmt.Printf("INFO signed %s request from %s using client version %s\n", frame.Mode, info.ConnectionId, info.ClientVersion)
	fmt.Printf("INFO info %+v\n", info)
	return HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("{\"server\":\"2.3.2\",\"state\":\"%s\",\"context\":\"%s\"}", frame.State, signature),
	}
}

func lambdaHandler(ctx context.Context, rawEvent json.RawMessage) (HTTPResponse, error) {
	fmt.Println(string(rawEvent))
	event := Event{}
	_ = json.Unmarshal(rawEvent, &event)

	headers := map[string]string{
		"Strict-Transport-Security": "max-age=31536000; includeSubdomains; preload", // 1 year
	}

	response := HTTPResponse{
		StatusCode: 200,
		Headers:    headers,
	}

	switch event.Context.EventType {
	case EventTypeConnect, EventTypeDisconnect:
		// Do nothing
	case EventTypeMessage:
		response = processMessage(ctx, event)
		response.Headers = headers
		return response, nil
	default:
		response.StatusCode = http.StatusNotFound
	}

	return response, nil
}
