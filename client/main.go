package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

const UsageInfo = `Version: 1.0.2
Usage:
  aada -configure

When configure completes, it will list what Azure AD roles/groups you have and what profiles
they have been installed into.  You should see something ilke this:

AWS_012345678901_RoleName -> RoleName

You will find a profile in your ~/.aws/config file called "RoleName".  This profile will be 
configured to use aada to fetch credentials, meaning you can make any standard AWS call like
you already had credentials.  An easy starting point is:

aws --profile RoleName sts get-caller-identity

If the CLI needs to fetch credentials, a browser window will open to authenticate you.  The
credentials will be cached in ~/.aws/credentials for subsequent use.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(UsageInfo)
		return
	}

	err := internal()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func internal() error {
	// Generate a nonce for use later
	nonce := uuid.NewString()

	// Build the initial request
	frame := &Frame{
		Nonce:   nonce,
		Profile: os.Args[1],
		Mode:    "access",
	}

	switch strings.ToLower(os.Args[1]) {
	case "-uidemo":
		sp := spinner.New(spinner.CharSets[14], 100 * time.Millisecond)
		sp.Reverse()
		sp.Color("blue")
		sp.Suffix = " reticulating splines"
		sp.Start()
		time.Sleep(10 * time.Second)
		sp.Stop()
		color.Green("✓ reticulating splines")
		return nil
	case "-configure":
		frame.Mode = "configuration"
	default:
		err := lookupCache(frame)
		if err == nil {
			// We have cached credentials
			fmt.Println(frame.ToCredentialString())
			return nil
		}
	}

	// Start a websocket connection and send the nonce
	wss, err := startWebsocket()
	if err != nil {
		return errors.Wrap(err, "unable to initiate websocket")
	}
	if wss == nil {
		return errors.New("no websocket handle")
	}
	raw, err := json.Marshal(frame)
	if err != nil {
		return errors.Wrap(err, "unable to encode request")
	}
	err = wss.WriteMessage(websocket.TextMessage, raw)
	if err != nil {
		return errors.Wrap(err, "failed to write role")
	}
	mt, msg, err := wss.ReadMessage()
	if err != nil {
		return errors.Wrap(err, "failed to read state")
	}
	if mt != websocket.TextMessage {
		return errors.New("invalid message format")
	}
	err = json.Unmarshal(msg, frame)
	if err != nil {
		return errors.Wrap(err, "unable to unpack frame")
	}
	err = launchLogin(nonce, frame.State, frame.Mode == "configuration")
	if err != nil {
		return errors.Wrap(err, "failed to launch browser login")
	}
	mt, msg, err = wss.ReadMessage()
	if err != nil {
		return errors.Wrap(err, "failed to read response")
	}
	if mt != websocket.TextMessage {
		return errors.New("invalid message format")
	}
	err = json.Unmarshal(msg, frame)
	if err != nil {
		fmt.Println("Frame: ", string(msg))
		return errors.Wrap(err, "unable to unpack frame")
	}

	if frame.Mode == "access" {
		// We don't really care if the cache works, so ignore errors explicitly
		_ = cacheCredentials(frame)
		fmt.Println(frame.ToCredentialString())
		return nil
	}

	// We're doing configuration, so we should have a list of profiles to configure
	if len(frame.ProfileList) == 0 {
		fmt.Println("no profiles were found for your account")
		return nil
	}

	return setupProfiles(frame.ProfileList)
}

func startWebsocket() (*websocket.Conn, error) {
	wss, _, err := websocket.DefaultDialer.Dial("wss://wss.aabg.io", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to dial remote websocket")
	}
	return wss, nil
}

//const authUrl = "https://login.microsoftonline.com/f3211d0e-125b-42c3-86db-322b19a65a22/oauth2/v2.0/authorize"
const authUrl = "https://login.microsoftonline.com/e0793d39-0939-496d-b129-198edd916feb/oauth2/v2.0/authorize"

func launchLogin(nonce string, state string, requireConsent bool) error {
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
	//if requireConsent {
	//	rqv.Set("prompt", "consent")
	//}
	return browser.OpenURL(authUrl + "?" + rqv.Encode())
}
