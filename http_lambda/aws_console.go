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
						DisplayName: userGroup.FriendlyName,
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

	err = cacheServicePrincipalGroups(creds)
	if err != nil {
		buildFailureResponse("failed to fetch principals: " + err.Error())
	}

	userGroups, err := getUserGroups(creds)
	if err != nil {
		return buildFailureResponse("failed to get user groups: " + err.Error()), nil
	}

	content := make([]ConsoleLink, 0)

	for _, userGroup := range userGroups {
		for _, app := range azureAppDefCache {
			for _, assignment := range app.Assignments {
				if assignment.Id == userGroup.Id {
					link := fmt.Sprintf(acpLoginUrl, app.DisplayName, app.AppId)
					content = append(content, ConsoleLink{
						Url: link,
						DisplayName: userGroup.FriendlyName,
						Account: userGroup.AccountId,
					})
				}
				break
			}
		}
	}

	return buildConsolePage(content), nil
}

func init() {
	_ = cbor.Unmarshal(AzureAppsData, &AzureApps)
}
