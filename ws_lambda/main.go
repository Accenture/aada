package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"net/http"
)

func main() {
	lambda.Start(lambdaHandler)
}

type Frame struct {
	Nonce           string `json:"nonce"`
	Profile         string `json:"profile"`
	State           string `json:"state"`
	Mode            string `json:"mode"`
	Connection      string `json:"-"`
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

	return HTTPResponse{
		StatusCode:      http.StatusOK,
		Body:            fmt.Sprintf("{\"state\":\"%s\"}", frame.State),
	}
}

func lambdaHandler(ctx context.Context, rawEvent json.RawMessage) (HTTPResponse, error) {
	fmt.Println(string(rawEvent))
	event := Event{}
	_ = json.Unmarshal(rawEvent, &event)

	switch event.Context.EventType {
	case EventTypeConnect:
		return HTTPResponse{StatusCode: 200}, nil
	case EventTypeDisconnect:
		return HTTPResponse{StatusCode: 200}, nil
	case EventTypeMessage:
		return processMessage(ctx, event), nil
	default:
		return HTTPResponse{
			StatusCode: 404,
		}, nil
	}
}
