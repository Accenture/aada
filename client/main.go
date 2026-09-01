package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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

func startWebsocket() (*websocket.Conn, error) {
	wss, _, err := websocket.DefaultDialer.Dial("wss://wss.aabg.io", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to dial remote websocket")
	}
	return wss, nil
}

const authUrl = "https://login.microsoftonline.com/e0793d39-0939-496d-b129-198edd916feb/oauth2/v2.0/authorize"
