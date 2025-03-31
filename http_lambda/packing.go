package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// AWS-0123456789012-Role_Name_Here
func unpackGroupName(profile string) (string, string, error) {
	separator := "UNKNOWN"
	if strings.HasPrefix(profile, "AWS-") {
		separator = "-"
	}
	if strings.HasPrefix(profile, "AWS_") {
		separator = "_"
	}
	parts := strings.SplitN(profile, separator, 3) // note this only splits into up to three pieces
	if len(parts) != 3 {
		return "", "", errors.New("^invalid role structure " + profile)
	}
	if parts[0] != "AWS" {
		return "", "", errors.New("invalid role name")
	}
	// verify the account number is fully numeric
	if match, _ := regexp.MatchString("\\d+", parts[1]); !match {
		return "", "", errors.New("role does not include valid account number")
	}
	return parts[1], parts[2], nil
}

func extractUpn(jwt string) string {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		fmt.Println("TOKEN FAILURE", jwt)
		return "unknown"
	}
	part, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("DECODE FAILURE", jwt)
		return "unknown"
	}
	structure := struct {
		Upn string `json:"upn"`
	}{}
	err = json.Unmarshal(part, &structure)
	if err != nil {
		fmt.Println("MARSHAL FAILURE", jwt)
		return "unknown"
	}
	return structure.Upn
}
