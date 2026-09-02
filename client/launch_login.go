package main

import (
	"fmt"
	"net/url"

	"github.com/pkg/browser"
)

func launchLogin(nonce string, state string, requireConsent bool, cli bool, loginHint string, prompt string) error {
	rqv := url.Values{}
	rqv.Set("nonce", nonce)
	rqv.Set("state", state)
	rqv.Set("client_id", "dbf2de86-2e04-4086-bc86-bbc8b47076d5")
	//rqv.Set("response_type", "code id_token")
	//rqv.Set("response_mode", "form_post")
	rqv.Set("response_type", "code")
	rqv.Set("response_mode", "query")
	rqv.Set("scope", "openid profile email")
	rqv.Set("redirect_uri", "https://aabg.io/authenticator")
	if loginHint != "" {
		rqv.Set("login_hint", loginHint)
	}
	if prompt != "" {
		rqv.Set("prompt", prompt)
	}
	//if requireConsent {
	//	rqv.Set("prompt", "consent")
	//}
	if cli {
		fmt.Println("Please open the following URL in your browser:")
		fmt.Println(authUrl + "?" + rqv.Encode())
		return nil
	} else {
		return browser.OpenURL(authUrl + "?" + rqv.Encode())
	}
}
