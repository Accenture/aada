package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/pkg/errors"
)

func assumeRole(ctx context.Context, upn string, accountId string, roleName string, duration int) (*types.Credentials, error) {
	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, roleName)

	svc := sts.NewFromConfig(awsConfig)
	assumeRoleInput := &sts.AssumeRoleInput{
		RoleArn:         &arn,
		RoleSessionName: &upn,
	}
	// If a duration has been specified, ask for it explicitly.  Zero means leave it off and let the default occur.
	if duration > 0 {
		assumeRoleInput.DurationSeconds = aws.Int32(int32(duration))
	}
	aro, err := svc.AssumeRole(ctx, assumeRoleInput)
	if err != nil {
		return nil, errors.Wrap(err, "role assumption failed")
	}

	fmt.Printf("AUDIT %s assumed %s in %s with key %s\n", upn, roleName, accountId, *aro.Credentials.AccessKeyId)
	return aro.Credentials, nil
}

func assumeRoleWithWebIdentity(ctx context.Context, upn string, accountId string, roleName string, token string) (*types.Credentials, error) {
	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, roleName)

	svc := sts.NewFromConfig(awsConfig)
	aro, err := svc.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          &arn,
		RoleSessionName:  &upn,
		WebIdentityToken: &token,
	})
	if err != nil {
		return nil, errors.Wrap(err, "role assumption failed")
	}
	fmt.Printf("AUDIT %s assumed %s in %s with key %s using token\n", upn, roleName, accountId, *aro.Credentials.AccessKeyId)
	return aro.Credentials, nil
}
