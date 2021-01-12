package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	username := flag.String("u", "", "specify the username")
	duration := flag.Duration("d", 1*time.Hour, "duration of assumed credentials")
	trace := flag.String("trace", "", "if specified, traces all activity into this file")
	version := flag.Bool("version", false, "display version information and exit")
	flag.Parse()

	if *version {
		fmt.Println("aaca 0.1.4")
		fmt.Println("  This tool authenticates with the Accenture federation service to obtain a SAML token")
		fmt.Println("  that is exchanged for AWS credentials.  Those credentials are written into the AWS")
		fmt.Println("  CLI/SDK credentials file for use by the CLI or other applications that use the SDK.")
		os.Exit(0)
	}

	in := bufio.NewReader(os.Stdin)

	var un string
	fmt.Print("Username: ")
	if *username == "" {
		var err error
		un, err = in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		un = strings.Trim(un, "\r\n")
	} else {
		un = *username
		fmt.Println(*username)
	}
	fmt.Print("Password: ")
	pw, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("<redacted>")

	fmt.Print("Symantec VIP: ")
	vip, err := in.ReadString('\n')
	vip = strings.Trim(vip, "\r\n")

	fmt.Println("0/5 preparing")
	awsConfig, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Panic(err)
	}
	awsConfig.Region = endpoints.UsEast1RegionID

	a := &Authenticator{
		Username:  un,
		Password:  string(pw),
		VIPCode:   vip,
		AWSConfig: awsConfig,
	}

	if *trace != "" {
		Trace, err = os.Create(*trace)
		if err != nil {
			log.Fatal(err)
		}
		defer Trace.Close()
		fmt.Printf("!! tracing enabled to %s\n", *trace)
	}

	fmt.Println("1/5 getting new session")
	a.getSubmissionFormUrl()
	fmt.Println("2/5 authenticating")
	a.submitCredentials()
	fmt.Println("3/5 fetching SAML token")
	a.submitVipCode()

	roles := a.getSAMLRoles()
	if len(roles) == 0 {
		log.Fatal("failed to unpack roles from the SAML response")
	}
	fmt.Println("You may assume one of the following roles:")
	for i, r := range roles {
		cutName := strings.Split(r, ",")
		fmt.Printf("  %2d: %s\n", i, cutName[1])
	}
	fmt.Print("Role Number: ")
	role, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	role = strings.Trim(role, "\r\n")
	roleId, err := strconv.Atoi(role)
	if err != nil {
		log.Fatal(err)
	}
	if roleId < 0 || roleId >= len(roles) {
		log.Fatal("invalid role id requested")
	}
	fmt.Println("4/5 exchanging SAML token for role credentials")
	parts := strings.Split(roles[roleId], ",")
	a.assumeRole(parts[0], parts[1], *duration)
	profile := parts[1]
	profile = profile[strings.LastIndex(profile, "/")+1:] // Get just the profile name after role/
	fmt.Printf("5/5 installing access key %s into profile %s\n", a.accessKey, profile)
	err = updateProfile(profile, a.accessKey, a.secretKey, a.sessionToken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("complete - your access expires in %s\n", time.Until(a.expiration).Round(1*time.Second))
}
