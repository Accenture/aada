package main

import (
	"bytes"
	"github.com/go-ini/ini"
	"io/ioutil"
	"os"
	"path"
)

func updateProfile(profile string, key string, secret string, token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
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

	// Ensure there is a profile section in the config with the proper name
	section := f.Section(profile)
	// Only update the output if it doesn't already exist
	if !section.HasKey("output") {
		_, err = section.NewKey("output", "json")
		if err != nil {
			return err
		}
	}
	// Only update the region if it doesn't already exist
	if !section.HasKey("region") {
		_, err = section.NewKey("region", "us-east-1")
		if err != nil {
			return err
		}
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

	f, err = ini.Load(credsFile)
	if err != nil {
		return err
	}

	section = f.Section(profile)
	_, err = section.NewKey("aws_access_key_id", key)
	if err != nil {
		return err
	}
	_, err = section.NewKey("aws_secret_access_key", secret)
	if err != nil {
		return err
	}
	_, err = section.NewKey("aws_session_token", token)
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
