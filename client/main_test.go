package main

import (
	"github.com/google/uuid"
	"testing"
)

func TestLaunchWithContext(t *testing.T) {
	nonce := uuid.NewString()
	packedState := "ogVwR29GNHNlVnlJQU1DRXhnPQZYRjBEAiByOnYFTNZwycuraJR3sDp9URiMQpDKnmSeDVzrseyhQwIgb1/PS6q2Z5SvCd7Sp/WxF9NW2Z2IIabJ0uwa6mEdnwE="
	err := launchLogin(nonce, packedState, false)
	if err != nil {
		t.Fatal(err)
	}
}
