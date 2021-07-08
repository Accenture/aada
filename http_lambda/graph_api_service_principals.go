package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

const servicePrincipalQuery = "https://graph.microsoft.com/v1.0/servicePrincipals?$count=true&$filter=tags%2Fany%28t%3At%20eq%20%27AWS%27%29&$select=id,appId,displayName,tags"
const servicePrincipalQueryByAccount = "https://graph.microsoft.com/v1.0/servicePrincipals?$count=true&$filter=tags%2Fany%28t%3At%20eq%20%27AWSAccountNumber%3D__ACCOUNT__%27%29&$select=id,appId,displayName,tags"
const appRoleAssignmentsQuery = "https://graph.microsoft.com/v1.0/servicePrincipals/%s/appRoleAssignedTo"

type AzureAppDef struct {
	Id          string   `json:"id"`
	AppId       string   `json:"appId"`
	DisplayName string   `json:"displayName"`
	Tags        []string `json:"tags"`
	Assignments []AzureAppRoleAssignment
}

type AzureAppRoleAssignment struct {
	Id                   string `json:"id"`
	AppRoleId            string `json:"appRoleId"`
	PrincipalDisplayName string `json:"principalDisplayName"`
	PrincipalId          string `json:"principalId"`
}

var azureAppDefCacheLock sync.Mutex
var azureAppDefCache map[string][]AzureAppDef

func fetchAssignedGroupForAWSAccount(creds *Credentials, accountNumber string) []AzureAppDef {
	// Cache check
	azureAppDefCacheLock.Lock()
	defs, ok := azureAppDefCache[accountNumber]
	azureAppDefCacheLock.Unlock()
	if ok {
		return defs
	}

	servicePrincipals := make([]AzureAppDef, 0)
	servicePrincipalsJson, errs := loadGraphResultSet(creds, strings.ReplaceAll(servicePrincipalQueryByAccount, "__ACCOUNT__", accountNumber))

	for servicePrincipalJson := range servicePrincipalsJson {
		servicePrincipal := AzureAppDef{}
		err := json.Unmarshal(servicePrincipalJson, &servicePrincipal)
		if err != nil {
			fmt.Println("ERROR failed to unmarshal service principal")
			continue
		}

		appRoles, subErr := loadGraphResultSet(creds, fmt.Sprintf(appRoleAssignmentsQuery, servicePrincipal.Id))

		for appRoleJson := range appRoles {
			appRoleAssignment := AzureAppRoleAssignment{}
			err := json.Unmarshal(appRoleJson, &appRoleAssignment)
			if err != nil {
				fmt.Println("ERROR failed to unmarshal app role assignment")
				continue
			}
			servicePrincipal.Assignments = append(servicePrincipal.Assignments, appRoleAssignment)
		}

		servicePrincipals = append(servicePrincipals, servicePrincipal)

		for err := range subErr {
			fmt.Println("ERROR", err.Error())
		}
	}

	for err := range errs {
		fmt.Println("ERROR", err.Error())
	}

	azureAppDefCacheLock.Lock()
	azureAppDefCache[accountNumber] = servicePrincipals
	azureAppDefCacheLock.Unlock()

	return servicePrincipals
}

func init() {
	azureAppDefCache = make(map[string][]AzureAppDef)
}
