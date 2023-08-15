package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"net/http"
	"time"
)

type ActiveState struct {
	State      string
	Nonce      string
	Profile    string
	Mode       string
	Connection string
	Gateway    string
	Region     string
}

func processOIDCRequest(ctx context.Context, state string, code string, idToken string) (Response, error) {
	var activeState *ActiveState

	si := &SignedInformation{}
	err := si.DecodeFromString(state)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return buildFailureResponse("failed to unpack state"), nil
	}
	err = si.Validate(ctx)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return buildFailureResponse("failed to validate state"), nil
	}

	// At this point, passed in state is valid and verified, proceed to trust it
	activeState = &ActiveState{
		Profile:    si.Information.ProfileName,
		Connection: si.Information.ConnectionId,
		Region:     si.Information.AWSRegion,
		Gateway:    si.Information.ApiId,
	}
	switch si.Information.ConnectMode {
	case ModeAccess:
		activeState.Mode = "access"
	case ModeConfiguration:
		activeState.Mode = "configuration"
	}

	customMessage := "login successful"

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
		//fmt.Println("DEBUG configuration token", accessToken)
		profiles, err := getUserProfiles(accessToken)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return buildFailureResponse("failed to query group membership"), nil
		}
		ar.ProfileList = profiles
	case "access":
		//fmt.Println("DEBUG access request for profile", activeState.Profile)
		ok, err := checkUserInsideGroup(accessToken, activeState.Profile)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return buildFailureResponse("failed to query group membership"), nil
		}
		if !ok {
			err = sendMessageToClient(ctx, activeState.Region, activeState.Gateway, activeState.Connection, NotAllowed)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				return buildFailureResponse("failed to validate group membership"), nil
			}
		}
		accountId, groupName, err := unpackGroupName(activeState.Profile)
		if err != nil {
			fmt.Println("ERROR", err.Error())
			_ = sendMessageToClient(ctx, activeState.Region, activeState.Gateway, activeState.Connection, InvalidGroupName)
			return buildFailureResponse("failed to unpack group name"), nil
		}
		var tok *types.Credentials
		upn := extractUpn(accessToken.AccessToken)
		if len(idToken) > 0 {
			tok, err = assumeRoleWithWebIdentity(ctx, upn, accountId, groupName, idToken)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				_ = sendMessageToClient(ctx, activeState.Region, activeState.Gateway, activeState.Connection, RoleAssumptionFailure)
				return buildFailureResponse("failed to fetch credentials"), nil
			}
		} else {
			tok, err = assumeRole(ctx, upn, accountId, groupName)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				_ = sendMessageToClient(ctx, activeState.Region, activeState.Gateway, activeState.Connection, RoleAssumptionFailure)
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
	_ = sendMessageToClient(ctx, activeState.Region, activeState.Gateway, activeState.Connection, string(msg))

	return buildSuccessResponse(customMessage), nil
}
