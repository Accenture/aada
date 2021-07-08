package main

import (
	"fmt"
	"testing"
)

func TestGroupThingy(t *testing.T) {
	creds := &Credentials{
		TokenType:   "Bearer",
		ExpiresIn:   0,
		AccessToken: "eyJ0eXAiOiJKV1QiLCJub25jZSI6IlRyejBKN1JwancyWHI1b0R2VVZ5ZVYzZkVybHdxMWt2bFBaNDE0RXZHUlEiLCJhbGciOiJSUzI1NiIsIng1dCI6Im5PbzNaRHJPRFhFSzFqS1doWHNsSFJfS1hFZyIsImtpZCI6Im5PbzNaRHJPRFhFSzFqS1doWHNsSFJfS1hFZyJ9.eyJhdWQiOiIwMDAwMDAwMy0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9lMDc5M2QzOS0wOTM5LTQ5NmQtYjEyOS0xOThlZGQ5MTZmZWIvIiwiaWF0IjoxNjI1NjkzOTI0LCJuYmYiOjE2MjU2OTM5MjQsImV4cCI6MTYyNTY5NzgyNCwiYWNjdCI6MCwiYWNyIjoiMSIsImFjcnMiOlsidXJuOnVzZXI6cmVnaXN0ZXJzZWN1cml0eWluZm8iLCJ1cm46bWljcm9zb2Z0OnJlcTEiLCJ1cm46bWljcm9zb2Z0OnJlcTIiLCJ1cm46bWljcm9zb2Z0OnJlcTMiLCJjMSIsImMyIiwiYzMiLCJjNCIsImM1IiwiYzYiLCJjNyIsImM4IiwiYzkiLCJjMTAiLCJjMTEiLCJjMTIiLCJjMTMiLCJjMTQiLCJjMTUiLCJjMTYiLCJjMTciLCJjMTgiLCJjMTkiLCJjMjAiLCJjMjEiLCJjMjIiLCJjMjMiLCJjMjQiLCJjMjUiXSwiYWlvIjoiQVZRQXEvOFRBQUFBYXJYdTRRbGYxdTBPU1RrczhQdHlZVDVTYmNWRG9MRDZpSnFUaktWbmM4WGtqWFBKeENJMHVuVXhKK3RRSGFrejU1NDVxR1kxYWtNbU03QkV5R28xd08wTk1lL01oTUhHVWRtWnFTK0VjbGc9IiwiYW1yIjpbInB3ZCIsIm1mYSJdLCJhcHBfZGlzcGxheW5hbWUiOiJHcmFwaCBleHBsb3JlciAob2ZmaWNpYWwgc2l0ZSkiLCJhcHBpZCI6ImRlOGJjOGI1LWQ5ZjktNDhiMS1hOGFkLWI3NDhkYTcyNTA2NCIsImFwcGlkYWNyIjoiMCIsImNvbnRyb2xzIjpbImFwcF9yZXMiXSwiY29udHJvbHNfYXVkcyI6WyIwMDAwMDAwMy0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCIwMDAwMDAwMy0wMDAwLTBmZjEtY2UwMC0wMDAwMDAwMDAwMDAiXSwiZGV2aWNlaWQiOiIwOWJiNTg1Ny0xMTYzLTQyZWUtOGQ4OS04YjczZTk3ZGU0MDAiLCJmYW1pbHlfbmFtZSI6IkhpbGwiLCJnaXZlbl9uYW1lIjoiRXJpYyIsImlkdHlwIjoidXNlciIsImluX2NvcnAiOiJ0cnVlIiwiaXBhZGRyIjoiMTA0LjYyLjEzOS45NCIsIm5hbWUiOiJIaWxsLCBFcmljIiwib2lkIjoiZTMxY2Y0YjctNzI1Zi00ZWQ5LWE3ZjYtMzcxYTUyMzVkMTlhIiwib25wcmVtX3NpZCI6IlMtMS01LTIxLTMyOTA2ODE1Mi0xNDU0NDcxMTY1LTE0MTcwMDEzMzMtNzM3NjE4MiIsInBsYXRmIjoiNSIsInB1aWQiOiIxMDAzMDAwMEFEODEzOUFFIiwicmgiOiIwLkFYc0FPVDE1NERrSmJVbXhLUm1PM1pGdjY3WElpOTc1MmJGSXFLMjNTTnB5VUdSN0FLSS4iLCJzY3AiOiJBcHBsaWNhdGlvbi5SZWFkLkFsbCBBdWRpdExvZy5SZWFkLkFsbCBDYWxlbmRhcnMuUmVhZFdyaXRlIENvbnRhY3RzLlJlYWRXcml0ZSBEaXJlY3RvcnkuQWNjZXNzQXNVc2VyLkFsbCBEaXJlY3RvcnkuUmVhZFdyaXRlLkFsbCBGaWxlcy5SZWFkV3JpdGUuQWxsIEdyb3VwLlJlYWRXcml0ZS5BbGwgSWRlbnRpdHlSaXNrRXZlbnQuUmVhZC5BbGwgTWFpbC5SZWFkV3JpdGUgTWVtYmVyLlJlYWQuSGlkZGVuIE5vdGVzLlJlYWRXcml0ZS5BbGwgb3BlbmlkIFBlb3BsZS5SZWFkIFBvbGljeS5SZWFkV3JpdGUuQ29uZGl0aW9uYWxBY2Nlc3MgcHJvZmlsZSBTaXRlcy5SZWFkV3JpdGUuQWxsIFRhc2tzLlJlYWRXcml0ZSBVc2VyLlJlYWQgVXNlci5SZWFkLkFsbCBVc2VyLlJlYWRCYXNpYy5BbGwgVXNlci5SZWFkV3JpdGUgVXNlci5SZWFkV3JpdGUuQWxsIGVtYWlsIiwic2lnbmluX3N0YXRlIjpbImR2Y19tbmdkIiwiZHZjX2NtcCIsImttc2kiXSwic3ViIjoiY0Nfd2pYa0pFLXZxRkhQek1BQ05NdnNDRXRTQ09mX3hpLVc0aHZCdU5oQSIsInRlbmFudF9yZWdpb25fc2NvcGUiOiJXVyIsInRpZCI6ImUwNzkzZDM5LTA5MzktNDk2ZC1iMTI5LTE5OGVkZDkxNmZlYiIsInVuaXF1ZV9uYW1lIjoiZXJpYy5oaWxsQGFjY2VudHVyZS5jb20iLCJ1cG4iOiJlcmljLmhpbGxAYWNjZW50dXJlLmNvbSIsInV0aSI6IjlJOHAyX3N0ZWstenJuSkNqMWpHQVEiLCJ2ZXIiOiIxLjAiLCJ3aWRzIjpbImI3OWZiZjRkLTNlZjktNDY4OS04MTQzLTc2YjE5NGU4NTUwOSJdLCJ4bXNfc3QiOnsic3ViIjoiQzRYWGpBY0JtT0VPd21Tdm43VjdMZ3g5dlZFYXNBUWRGWlQyR0JBd1BsayJ9LCJ4bXNfdGNkdCI6MTM5NjA0OTM2NH0.nm_IP1p6tQiTsg4f1OiZj7n5cCzBig9R5ey20_oFYNSLyl2hu32_MuxIfqYVOTuRunsSJzA4AihI7gOlqAqkFHOVJvqAOh3OXO0OhMwCKXo6zpyXKHWwnF4lJBGVHMHrOi3zfHsrBUzUqsPz4_plaQx21vv_gC4KB5n8_1V3k9eeNnheElj_XDXVZo096gxT78J18Y8rxqZS74X4h716kJu-ZWaRIbgibd7Yyo-HPJuCBf4GEDI2NG-YtzcM-GbCx4GU_bklQWbSj8CYBv3NRZkp34u-bGuuY1iC_D6ZPS3w-u0SfluFAF4NqFcbKPVWyENANBlMuKzeb1bAqg0Gvg",
	}

	content := make([]ConsoleLink, 0)

	// Get a list of user groups
	userGroups, err := getUserGroups(creds)
	if err != nil {
		t.Fatal(err)
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
			t.Logf("principal %s", servicePrincipal.DisplayName)
			link := ConsoleLink{
				Account: account,
				Url: fmt.Sprintf(acpLoginUrl, servicePrincipal.DisplayName, servicePrincipal.AppId),
			}
			for _, userGroup := range groupList {
				for _, assignment := range servicePrincipal.Assignments {
					t.Logf("  matching assignment %+v", assignment)
					if assignment.PrincipalId == userGroup.Id {
						link.DisplayNames = append(link.DisplayNames, userGroup.FriendlyName)
					}
					if userGroup.FriendlyName == "Wade_Dev_Admin" {
						link.DisplayNames = append(link.DisplayNames, "Wade_Other_Admin")
					}
				}
			}
			if len(link.DisplayNames) > 0 {
				content = append(content, link)
			}
		}
	}

	t.Logf("built %+v", content)
}
