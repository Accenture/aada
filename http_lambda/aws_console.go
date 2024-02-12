package main

import (
	_ "embed"
	"fmt"
	"sync"
	"time"
)

const acpLoginUrl = "https://myapps.microsoft.com/signin/%s/%s?tenantId=e0793d39-0939-496d-b129-198edd916feb"

func buildAWSConsoleDisplay(code string) (Response, error) {
	creds, err := getAccessTokenFromCode(code)
	if err != nil {
		return buildFailureResponse("failed to convert code: " + err.Error()), nil
	}

	upn := extractUpn(creds.AccessToken)

	cachedContent := fetchCachedConsoleContent(upn)
	if cachedContent != nil {
		fmt.Printf("AUDIT %s accessed the console with %d entries via cached data\n", upn, len(cachedContent))

		return buildConsolePage(cachedContent), nil
	}

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

	workers := sync.WaitGroup{}
	contentLock := sync.Mutex{}
	content := make([]ConsoleLink, 0)

	for account, groupList := range groupsByAccount {
		//fmt.Println("DEBUG inside group loop", account)
		workers.Add(1)
		account := account
		groupList := groupList
		go func() {
			//fmt.Println("DEBUG fetchAssignedGroupForAWSAccount", account)
			for _, servicePrincipal := range fetchAssignedGroupForAWSAccount(creds, account) {
				link := ConsoleLink{
					Account: account,
					Url:     fmt.Sprintf(acpLoginUrl, servicePrincipal.DisplayName, servicePrincipal.AppId),
				}
				for _, userGroup := range groupList {
					for _, assignment := range servicePrincipal.Assignments {
						if assignment.PrincipalId == userGroup.Id {
							link.DisplayNames = append(link.DisplayNames, userGroup.FriendlyName)
						}
					}
				}
				if len(link.DisplayNames) > 0 {
					contentLock.Lock()
					content = append(content, link)
					contentLock.Unlock()
				}
			}
			workers.Done()
		}()
	}
	workers.Wait()

	cacheConsoleContent(upn, content)

	fmt.Printf("AUDIT %s accessed the console with %d entries\n", upn, len(content))

	return buildConsolePage(content), nil
}

type cachedEntry struct {
	expires time.Time
	content []ConsoleLink
}

var cache map[string]cachedEntry = make(map[string]cachedEntry)

func cacheConsoleContent(upn string, content []ConsoleLink) {
	cache[upn] = cachedEntry{
		expires: time.Now().Add(5 * time.Minute),
		content: content,
	}
}

func fetchCachedConsoleContent(upn string) []ConsoleLink {
	hit, ok := cache[upn]
	if ok {
		if hit.expires.After(time.Now()) {
			return hit.content
		}
	}
	return nil
}
