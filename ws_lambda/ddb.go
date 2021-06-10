package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"os"
	"strconv"
	"time"
)

var tableName string

func init() {
	var ok bool
	tableName, ok = os.LookupEnv("TABLE_NAME")
	if !ok {
		fmt.Print("no table name was specified in TABLE_NAME\n")
	}
}

func (frame *Frame) Persist(ctx context.Context) error {
	ddb := dynamodb.NewFromConfig(awsConfig)

	// Requests automatically expire an hour into the future
	expiration := strconv.FormatInt(time.Now().Add(1 * time.Hour).Unix(), 10)

	_, err := ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item: map[string]types.AttributeValue{
			"state": &types.AttributeValueMemberS{Value: frame.State},
			"nonce": &types.AttributeValueMemberS{Value: frame.Nonce},
			"profile": &types.AttributeValueMemberS{Value: frame.Profile},
			"mode": &types.AttributeValueMemberS{Value: frame.Mode},
			"connection": &types.AttributeValueMemberS{Value: frame.Connection},
			"expiration": &types.AttributeValueMemberN{Value: expiration},
		},
	})
	return err
}

func (frame *Frame) ToJson() (string, error) {
	raw, err := json.Marshal(frame)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
