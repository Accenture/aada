package main

import "testing"

func TestSomething(t *testing.T) {
	creds := &Credentials{
		TokenType:   "Bearer",
		ExpiresIn:   0,
		AccessToken: "yearight",
	}
	err := cacheServicePrincipalGroups(creds)
	if err != nil {
		t.Fatal(err)
	}
}
