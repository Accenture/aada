package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/pkg/errors"
)

func assumeRole(ctx context.Context, upn string, accountId string, roleName string) (*types.Credentials, error) {
	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, roleName)

	svc := sts.NewFromConfig(awsConfig)
	aro, err := svc.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &arn,
		RoleSessionName: &upn,
	})
	if err != nil {
		return nil, errors.Wrap(err, "role assumption failed")
	}
	return aro.Credentials, nil
}

func assumeRoleWithWebIdentity(ctx context.Context, upn string, accountId string, roleName string, token string) (*types.Credentials, error) {
	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, roleName)

	svc := sts.NewFromConfig(awsConfig)
	aro, err := svc.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
		RoleArn: &arn,
		RoleSessionName: &upn,
		WebIdentityToken: &token,

	})
	if err != nil {
		return nil, errors.Wrap(err, "role assumption failed")
	}
	return aro.Credentials, nil
}