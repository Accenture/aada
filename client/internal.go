package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
)

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
	cliMode := false
	horizon := time.Now()
	loginHint := ""

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
			} else if strings.HasPrefix(strings.ToLower(os.Args[i]), "-cli") {
				cliMode = true
			} else if strings.HasPrefix(strings.ToLower(os.Args[i]), "-login-hint") {
				loginHint = os.Args[i][len("-login-hint")+1:]
				matched, err := regexp.MatchString("[a-z0-9\\\\.]+@[a-z0-9]+\\.[a-z0-9]+", loginHint)
				if err != nil {
					fmt.Println("failed to parse login regexp")
					return nil
				}
				if !matched {
					loginHint = ""
					// silently continue without a login hint
				}
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
	err = launchLogin(nonce, frame.Context, frame.Mode == "configuration", cliMode, loginHint)
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
