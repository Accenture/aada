package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestNameUnpacking(t *testing.T) {
	raw, err := ioutil.ReadFile("/Users/eric.hill/Library/Application Support/JetBrains/GoLand2021.3/scratches/scratch.json")
	if err != nil {
		t.Fatal(err)
	}

	attr := struct {
		Count    int    `json:"@odata.count"`
		NextLink string `json:"@odata.nextLink"`
		Values   []struct {
			Id          string `json:"id"`
			DisplayName string `json:"displayName"` // This won't come in yet, but it's useful to have
		} `json:"value"`
	}{}
	err = json.Unmarshal(raw, &attr)
	if err != nil {
		t.Fatal("unable to unmarshal groups")
	}
	profiles := make(map[string]string)

	for _, value := range attr.Values {
		// Unpack a reasonable name and map it
		_, _, err := unpackGroupName(value.DisplayName)
		if err == nil {
			profiles[value.DisplayName] = value.DisplayName
		}
	}

	t.Logf("profiles: %+v", profiles)
}
