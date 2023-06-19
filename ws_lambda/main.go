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

type Frame struct {
	Nonce         string `json:"nonce"`
	Profile       string `json:"profile"`
	State         string `json:"state"`
	Mode          string `json:"mode"`
	ClientVersion string `json:"client_version"`
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
	err = frame.Persist(ctx) // Save the state to Dynamo so we can lookup the profile later
	if err != nil {
		fmt.Println("error saving state")
		return HTTPResponse{StatusCode: http.StatusInternalServerError}
	}

	// Package up the information the caller needs to carry into a signed structure
	info := &Information{
		ConnectionId:     event.Context.ConnectionId,
		ProfileName:      frame.Profile,
		ConnectMode:      ModeUnknown,
		ClientVersion:    frame.ClientVersion,
		ConnectionTarget: websocketUrl,
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
			Body:       fmt.Sprintf("{\"version\":\"2.3.2\",\"state\":\"%s\"}", frame.State),
		}
	}
	return HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("{\"version\":\"2.3.2\",\"state\":\"%s\",\"context\":\"%s\"}", frame.State, signature),
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
