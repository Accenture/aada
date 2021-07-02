package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

const servicePrincipalQuery = "https://graph.microsoft.com/v1.0/servicePrincipals?$count=true&$filter=tags%2Fany%28t%3At%20eq%20%27AWS%27%29&$select=id,appId,displayName,tags"
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
var azureAppDefCache []AzureAppDef

func cacheServicePrincipalGroups(creds *Credentials) error {
	// No need to cache if we already have a cache
	if len(azureAppDefCache) > 0 {
		return nil
	}

	fmt.Println("DEBUG caching service principal groups")

	servicePrincipals, errs := loadGraphResultSet(creds, servicePrincipalQuery)

	// Make workers to get group information concurrently
	for worker := 0; worker < 20; worker++ {
		go func() {
			for servicePrincipalJson := range servicePrincipals {
				servicePrincipal := AzureAppDef{}
				err := json.Unmarshal(servicePrincipalJson, &servicePrincipal)
				if err != nil {
					fmt.Println("ERROR failed to unmarshal service principal")
					continue
				}

				appRoles, subErr := loadGraphResultSet(creds, fmt.Sprintf(appRoleAssignmentsQuery, servicePrincipal.Id))

				// Process the roles concurrently
				go func() {
					for appRoleJson := range appRoles {
						appRoleAssignment := AzureAppRoleAssignment{}
						err := json.Unmarshal(appRoleJson, &appRoleAssignment)
						if err != nil {
							fmt.Println("ERROR failed to unmarshal app role assignment")
							continue
						}
						servicePrincipal.Assignments = append(servicePrincipal.Assignments, appRoleAssignment)
					}

					azureAppDefCacheLock.Lock()
					azureAppDefCache = append(azureAppDefCache, servicePrincipal)
					azureAppDefCacheLock.Unlock()
				}()

				for err := range subErr {
					fmt.Println("ERROR", err.Error())
				}
			}
		}()
	}

	for err := range errs {
		fmt.Println("ERROR", err.Error())
	}

	return nil
}

