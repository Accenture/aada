package main

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"strings"
)

func validateGroupClaim(accessToken string) error {
	req, err := http.NewRequest("POST", "https://graph.microsoft.com/v1.0/me/getMemberGroups", strings.NewReader("{\"securityEnabledOnly\":true}"))
	if err != nil {
		return errors.Wrap(err, "building group request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	raw, err := httputil.DumpRequest(req, true)
	if err != nil {
		return errors.Wrap(err, "dumping group request")
	}
	fmt.Println("!! REQ", base64.StdEncoding.EncodeToString(raw))

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "making group request")
	}

	raw, err = httputil.DumpResponse(rsp, true)
	if err != nil {
		return errors.Wrap(err, "dumping group response")
	}

	fmt.Println("!! RSP", base64.StdEncoding.EncodeToString(raw))

	return nil
}
