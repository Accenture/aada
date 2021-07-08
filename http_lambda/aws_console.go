package main

import (
	_ "embed"
	"fmt"
	"github.com/fxamacker/cbor"
	"strings"
)

type AzureApp struct {
	AppId string `cbor:"1,keyasint,omitempty"`
	Name  string `cbor:"2,keyasint,omitempty"`
}

const acpLoginUrl = "https://myapps.microsoft.com/signin/%s/%s?tenantId=e0793d39-0939-496d-b129-198edd916feb"

//go:embed acp_apps.cbor
var AzureAppsData []byte
var AzureApps []AzureApp

func buildAWSConsoleDisplay(code string) (Response, error) {
	creds, err := getAccessTokenFromCode(code)
	if err != nil {
		return buildFailureResponse("failed to convert code: " + err.Error()), nil
	}

	userGroups, err := getUserGroups(creds)
	if err != nil {
		return buildFailureResponse("failed to query groups: " + err.Error()), nil
	}

	fmt.Printf("received %d groups\n", len(userGroups))

	content := make([]ConsoleLink, 0)
	dedup := make(map[string]interface{})

	for _, userGroup := range userGroups {
		for _, azureApp := range AzureApps {
			if strings.Contains(azureApp.Name, userGroup.AccountId) {
				link := fmt.Sprintf(acpLoginUrl, azureApp.Name, azureApp.AppId)
				_, ok := dedup[link]
				if !ok {
					content = append(content, ConsoleLink{
						Url:         link,
						DisplayNames: []string{userGroup.FriendlyName},
						Account:     userGroup.AccountId,
					})
					dedup[link] = nil
				}
				break
			}
		}
	}

	fmt.Printf("built %d links", len(content))

	return buildConsolePage(content), nil
}

func buildAWSConsoleDisplay2(code string) (Response, error) {
	creds, err := getAccessTokenFromCode(code)
	if err != nil {
		return buildFailureResponse("failed to convert code: " + err.Error()), nil
	}

	content := make([]ConsoleLink, 0)

	// Get a list of user groups
	userGroups, err := getUserGroups(creds)
	if err != nil {
		return buildFailureResponse("failed to fetch user groups"), nil
	}

	groupsByAccount := make(map[string][]UserGroupInfo)
	for _, userGroup := range userGroups {
		curr, ok := groupsByAccount[userGroup.AccountId]
		if ok {
			curr = append(curr, userGroup)
			groupsByAccount[userGroup.AccountId] = curr
		} else {
			groupsByAccount[userGroup.AccountId] = []UserGroupInfo{userGroup}
		}
	}

	for account, groupList := range groupsByAccount {
		for _, servicePrincipal := range fetchAssignedGroupForAWSAccount(creds, account) {
			link := ConsoleLink{
				Account: account,
				Url: fmt.Sprintf(acpLoginUrl, servicePrincipal.DisplayName, servicePrincipal.AppId),
			}
			for _, userGroup := range groupList {
				for _, assignment := range servicePrincipal.Assignments {
					if assignment.PrincipalId == userGroup.Id {
						link.DisplayNames = append(link.DisplayNames, userGroup.FriendlyName)
					}
				}
			}
			if len(link.DisplayNames) > 0 {
				content = append(content, link)
			}
		}
	}

	return buildConsolePage(content), nil
}

func init() {
	_ = cbor.Unmarshal(AzureAppsData, &AzureApps)
}
