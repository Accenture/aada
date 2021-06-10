package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"os"
)

var tableName string

func init() {
	var ok bool
	tableName, ok = os.LookupEnv("TABLE_NAME")
	if !ok {
		fmt.Print("no table name was specified in TABLE_NAME\n")
	}
}

type ActiveState struct {
	State      string
	Nonce      string
	Profile    string
	Mode       string
	Connection string
}

func loadState(state string) (*ActiveState, error) {
	ddb := dynamodb.NewFromConfig(awsConfig)

	out, err := ddb.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"state": &types.AttributeValueMemberS{Value: state},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch state")
	}

	activeState := &ActiveState{}
	_ = attributevalue.Unmarshal(out.Item["nonce"], &activeState.Nonce)
	_ = attributevalue.Unmarshal(out.Item["profile"], &activeState.Profile)
	_ = attributevalue.Unmarshal(out.Item["mode"], &activeState.Mode)
	_ = attributevalue.Unmarshal(out.Item["connection"], &activeState.Connection)

	return activeState, nil
}
