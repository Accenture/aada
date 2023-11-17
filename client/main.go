package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

//go:embed version.info
var version string

const UsageInfo = `
  __     __    ___    __   
 / /\   / /\  | | \  / /\  
/_/--\ /_/--\ |_|_/ /_/--\ 

Usage: aada -configure [-horizon=...] [-duration=...]

When configure completes, it will list what Azure AD roles/groups you have and what profiles
they have been installed into.  You should see something like this:

+-------------------------------------------+---------------------------------------+
|         AZURE AD APPLICATION NAME         |         AWS SDK PROFILE NAME          |
+-------------------------------------------+---------------------------------------+
| AWS_012345678901_RoleName                 | RoleName                              |
+-------------------------------------------+---------------------------------------+
|            PROFILES INSTALLED             |                  1                    |
+-------------------------------------------+---------------------------------------+

You will find a profile in your ~/.aws/config file called "RoleName".  This profile will be 
configured to use AADA to fetch credentials, meaning you can make any standard AWS call like
you already had credentials.  An easy starting point is:

aws --profile RoleName sts get-caller-identity

If the CLI needs to fetch credentials, a browser window will open to authenticate you.  The
credentials will be cached in ~/.aws/credentials for subsequent use.

If you need AADA to ensure a credential is valid for a minimum amount of time, such as when
running automation that takes 15-30 minutes to run, you can use the -horizon switch in the 
~/.aws/config file:

[profile example]
credential-process = aada AWS_01234567890_example -horizon=30m

This will ensure there's at least 30 minutes of time left on the returned credential before
it expires, or will request new credentials from the provider if the credential expires within
that time period.  You can use seconds (90s), minutes (30m), or hours (4h) with this switch.

There is also a -duration switch that requests the credentials for a specified duration of
time.  This is passed to the AssumeRole API call when requesting credentials, and must be 
less than or equal to the maximum session duration specified in the IAM role configuration.`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Version:", version)
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
		Nonce:         nonce,
		Profile:       os.Args[1],
		Mode:          "access",
		ClientVersion: version,
		Duration:      3600, // One hour by default
	}

	useLongNameFormat := false
	horizon := time.Now()

	for i := 1; i < len(os.Args); i++ {
		switch strings.ToLower(os.Args[i]) {
		case "console", "-console", "--console":
			err := browser.OpenURL("https://aabg.io/awsconsole")
			if err != nil {
				fmt.Println("failed to open https://aabg.io/awsconsole")
			}
			return err
		case "update", "-upgrade", "--upgrade":
			err := browser.OpenURL("https://aabg.io/downloads")
			if err != nil {
				fmt.Println("failed to open https://aabg.io/downloads")
			}
			return err
		case "configure", "-configure", "--configure":
			frame.Mode = "configuration"
		case "-long-profile-names", "--long-profile-names":
			useLongNameFormat = true
		case "version", "-v", "-version", "--version":
			fmt.Println("aada version", version)
			return nil
		case "-h", "-?", "-help", "--help":
			fmt.Println("Version:", version)
			fmt.Println(UsageInfo)
			return nil
		default:
			if strings.HasPrefix(strings.ToLower(os.Args[i]), "-horizon") {
				t, err := parseSwitch("horizon", os.Args[i])
				if err != nil {
					fmt.Println("failed to parse horizon")
					return nil
				}
				horizon = horizon.Add(t)
			} else if strings.HasPrefix(strings.ToLower(os.Args[i]), "-duration") {
				t, err := parseSwitch("duration", os.Args[i])
				if err != nil {
					fmt.Println("failed to parse duration")
					return nil
				}
				frame.Duration = int(t.Seconds())
			} else if os.Args[i][0:1] == "-" {
				fmt.Println("Invalid switch:", os.Args[i])
				fmt.Println(UsageInfo)
				return nil
			}
		}
	}

	if frame.Mode == "access" {
		err := lookupCache(frame, horizon)
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
	err = launchLogin(nonce, frame.Context, frame.Mode == "configuration")
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

	return setupProfiles(useLongNameFormat, frame.ProfileList)
}

func startWebsocket() (*websocket.Conn, error) {
	wss, _, err := websocket.DefaultDialer.Dial("wss://wss.aabg.io", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to dial remote websocket")
	}
	return wss, nil
}

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
