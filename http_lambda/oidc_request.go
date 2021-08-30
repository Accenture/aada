package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"net/http"
	"time"
)

func processOIDCRequest(ctx context.Context, state string, code string, idToken string, wsurl string) (Response, error) {
	activeState, err := loadState(state)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return buildFailureResponse("failed to load state"), nil
	}
	accessToken, err := getAccessTokenFromCode(code)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return buildFailureResponse("failed to fetch access token"), nil
	}
	ar := &AssumptionResponse{
		Status: "allowed",
	}

	switch activeState.Mode {
	case "configuration":
		profiles, err := getUserProfiles(accessToken)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return buildFailureResponse("failed to query group membership"), nil
		}
		ar.ProfileList = profiles
	case "access":
		ok, err := checkUserInsideGroup(accessToken, activeState.Profile)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return buildFailureResponse("failed to query group membership"), nil
		}
		if !ok {
			err = sendMessageToClient(ctx, wsurl, activeState.Connection, NotAllowed)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				return buildFailureResponse("failed to validate group membership"), nil
			}
		}
		upn := extractUpn(accessToken.AccessToken)
		accountId, groupName, err := unpackGroupName(activeState.Profile)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			_ = sendMessageToClient(ctx, wsurl, activeState.Connection, InvalidGroupName)
			return buildFailureResponse("failed to unpack group name"), nil
		}
		var tok *types.Credentials
		if len(idToken) > 0 {
			tok, err = assumeRoleWithWebIdentity(ctx, upn, accountId, groupName, idToken)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				_ = sendMessageToClient(ctx, wsurl, activeState.Connection, RoleAssumptionFailure)
				return buildFailureResponse("failed to fetch credentials"), nil
			}
		} else {
			tok, err = assumeRole(ctx, upn, accountId, groupName)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				_ = sendMessageToClient(ctx, wsurl, activeState.Connection, RoleAssumptionFailure)
				return buildFailureResponse("failed to fetch credentials"), nil
			}
		}

		ar.Profile = activeState.Profile
		ar.Version = 1
		ar.AccessKeyId = *tok.AccessKeyId
		ar.SecretAccessKey = *tok.SecretAccessKey
		ar.SessionToken = *tok.SessionToken
		ar.Expiration = tok.Expiration.Format(time.RFC3339)
	}

	msg, err := json.Marshal(ar)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return Response{StatusCode: http.StatusInternalServerError}, nil
	}
	_ = sendMessageToClient(ctx, wsurl, activeState.Connection, string(msg))
	return buildSuccessResponse(), nil
}
