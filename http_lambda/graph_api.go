package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const alternateQuery = "https://graph.microsoft.com/v1.0/users/e31cf4b7-725f-4ed9-a7f6-371a5235d19a/getMemberGroups"
const memberQuery = "https://graph.microsoft.com/v1.0/me/transitiveMemberOf?$search=\"displayName:%s\"&$count=true"
const groupListQuery = "https://graph.microsoft.com/v1.0/me/transitiveMemberOf/microsoft.graph.group?$select=id,displayName"
const groupNameQuery = "https://graph.microsoft.com/v1.0/users/%s/memberOf?$select=id,displayName"
const groupNameQuery2 = "https://graph.microsoft.com/v1.0/groups/%s?$select=displayName"
const tokenUrl = "https://login.microsoftonline.com/e0793d39-0939-496d-b129-198edd916feb/oauth2/v2.0/token" // Prod
//const tokenUrl = "https://login.microsoftonline.com/f3211d0e-125b-42c3-86db-322b19a65a22/oauth2/v2.0/token" // Staging

func getAccessTokenFromCode(code string) (*Credentials, error) {
	rqv := url.Values{}
	rqv.Set("code", code)
	rqv.Set("redirect_uri", "https://aabg.io/authenticator")
	rqv.Set("grant_type", "authorization_code")
	rqv.Set("client_id", "dbf2de86-2e04-4086-bc86-bbc8b47076d5")
	rqv.Set("scope", "openid email")

	clientSecret, _ := os.LookupEnv("CLIENT_SECRET")
	rqv.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(rqv.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create token conversion request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert code to access token")
	}
	rb, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read token body")
	}
	_ = rsp.Body.Close()
	bt := &Credentials{}
	err = json.Unmarshal(rb, &bt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal token body")
	}
	if len(bt.AccessToken) == 0 {
		fmt.Println("token exchange error " + string(rb))
	}
	return bt, nil
}

func getGroupName(creds *Credentials, groupId string) (string, error) {
	fmt.Println("getting group name", groupId)

	req, err := http.NewRequest("GET", fmt.Sprintf(groupNameQuery, groupId), nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to build request")
	}
	req.Header.Add("Authorization", creds.TokenType+" "+creds.AccessToken)
	req.Header.Add("ConsistencyLevel", "eventual")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute group query")
	}
	raw, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read response body")
	}

	dn := struct {
		Name string `json:"displayName"`
	}{}
	err = json.Unmarshal(raw, &dn)
	if err != nil {
		return "", errors.Wrap(err, "failed to unpack response")
	}
	if len(dn.Name) == 0 {
		fmt.Println(string(raw))
		return "", errors.New("api call response error")
	}

	return dn.Name, nil
}

func getUserProfiles(creds *Credentials) (map[string]string, error) {
	fmt.Println("fetching user profiles")

	profiles := make(map[string]string)

	req, err := http.NewRequest("GET", groupListQuery, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build request")
	}

	for {
		req.Header.Add("Authorization", creds.TokenType+" "+creds.AccessToken)
		req.Header.Add("ConsistencyLevel", "eventual")

		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "failed to execute group query")
		}
		raw, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read response body")
		}

		fmt.Println("RAW " + string(raw))

		attr := struct {
			Count    int    `json:"@odata.count"`
			NextLink string `json:"@odata.nextLink"`
			Values   []struct {
				Id          string `json:"id"`
				DisplayName string `json:"displayName"` // This won't come in yet, but it's useful to have
			} `json:"value"`
		}{}
		err = json.Unmarshal(raw, &attr)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal groups")
		}

		for _, value := range attr.Values {
			// Unpack a reasonable name and map it
			_, groupName, err := unpackGroupName(value.DisplayName)
			if err == nil {
				profiles[groupName] = value.DisplayName
			}
		}

		if attr.NextLink == "" {
			break
		}

		req, err = http.NewRequest("GET", attr.NextLink, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build subsequent request")
		}
	}

	return profiles, nil
}

func checkUserInsideGroup(creds *Credentials, groupName string) (bool, error) {
	// Fast fail if the group name includes any injection-like characters
	if strings.ContainsAny(groupName, "'\"!$:%") {
		return false, nil
	}
	query := fmt.Sprintf(memberQuery, groupName)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to build request")
	}
	req.Header.Add("Authorization", creds.TokenType+" "+creds.AccessToken)
	req.Header.Add("ConsistencyLevel", "eventual")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "failed to execute group query")
	}
	raw, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return false, errors.Wrap(err, "unable to read response body")
	}
	attr := struct {
		Count int `json:"@odata.count"`
	}{}
	err = json.Unmarshal(raw, &attr)
	if err != nil {
		return false, errors.Wrap(err, "unable to unmarshal count")
	}
	return attr.Count >= 1, nil // They are in the group if the count is 1
}