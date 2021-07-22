package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

// Either returns populated credentials, or nil if they don't exist
func lookupCache(frame *Frame) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	awsPath := path.Join(home, ".aws")
	if _, err := os.Stat(awsPath); os.IsNotExist(err) {
		_ = os.Mkdir(awsPath, 0755)
	}

	credsPath := path.Join(home, ".aws", "credentials")
	var credsFile []byte
	credsFile, err = ioutil.ReadFile(credsPath)
	if err != nil {
		if os.IsNotExist(err) {
			credsFile = []byte{} // Empty file
		} else {
			return err
		}
	}

	f, err := ini.Load(credsFile)
	if err != nil {
		return err
	}

	section := f.Section(frame.Profile)
	if !section.HasKey("expiration_date") {
		return errors.New("no expiration date, assuming credentials are stale")
	}
	exp, err := section.GetKey("expiration_date")
	if err != nil {
		return err
	}
	expt, err := time.Parse(time.RFC3339, exp.String())
	if err != nil {
		return err
	}
	if expt.Before(time.Now()) {
		return errors.New("credentials expired")
	}
	// Cached credentials are good, flesh out the rest of the response
	frame.Version = 1
	frame.AccessKeyId = section.Key("aws_access_key_id").String()
	frame.SecretAccessKey = section.Key("aws_secret_access_key").String()
	frame.SessionToken = section.Key("aws_session_token").String()
	frame.Expiration = exp.String()

	return nil
}

func cacheCredentials(frame *Frame) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	awsPath := path.Join(home, ".aws")
	if _, err := os.Stat(awsPath); os.IsNotExist(err) {
		_ = os.Mkdir(awsPath, 0755)
	}

	credsPath := path.Join(home, ".aws", "credentials")
	var credsFile []byte
	credsFile, err = ioutil.ReadFile(credsPath)
	if err != nil {
		if os.IsNotExist(err) {
			credsFile = []byte{} // Empty file
		} else {
			return err
		}
	}

	f, err := ini.Load(credsFile)
	if err != nil {
		return err
	}

	section := f.Section(frame.Profile)
	_, err = section.NewKey("aws_access_key_id", frame.AccessKeyId)
	if err != nil {
		return err
	}
	_, err = section.NewKey("aws_secret_access_key", frame.SecretAccessKey)
	if err != nil {
		return err
	}
	_, err = section.NewKey("aws_session_token", frame.SessionToken)
	if err != nil {
		return err
	}
	_, err = section.NewKey("expiration_date", frame.Expiration)
	if err != nil {
		return err
	}
	credsData := &bytes.Buffer{}
	_, err = f.WriteTo(credsData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(credsPath, credsData.Bytes(), 0660)
	if err != nil {
		return err
	}

	return nil
}

func setupProfiles(profiles map[string]string) error {
	fmt.Println("setting up profiles")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	awsPath := path.Join(home, ".aws")
	if _, err := os.Stat(awsPath); os.IsNotExist(err) {
		_ = os.Mkdir(awsPath, 0755)
	}

	configPath := path.Join(home, ".aws", "config")
	var configFile []byte
	configFile, err = ioutil.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			configFile = []byte{} // Empty file
		} else {
			return err
		}
	}

	f, err := ini.Load(configFile)
	if err != nil {
		return err
	}

	createdRoles := make(map[string]int)

	for profile, role := range profiles {
		n, ok := createdRoles[profile]
		if ok {
			// We already had a role by this name, so we increment
			n++
			createdRoles[profile] = n
			profile = profile + strconv.Itoa(n)
		} else {
			createdRoles[profile] = 1
		}

		section := f.Section("profile " + profile)
		// Clear out the existing keys
		for _, existingKey := range section.Keys() {
			section.DeleteKey(existingKey.Name())
		}
		_, err := section.NewKey("credential_process", "aada "+role)
		if err != nil {
			return err
		}

		fmt.Println("    ", role, "->", profile)
	}
	configData := &bytes.Buffer{}
	_, err = f.WriteTo(configData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, configData.Bytes(), 0660)
	if err != nil {
		return err
	}

	return nil
}
