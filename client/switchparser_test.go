package main

import (
	"testing"
	"time"
)

func TestSeconds(t *testing.T) {
	d, err := parseSwitch("duration", "-duration=5s")
	if err != nil {
		t.Fatal(err)
	}
	if d != (5 * time.Second) {
		t.Fatal("invalid duration parsed")
	}
}

func TestMinutes(t *testing.T) {
	d, err := parseSwitch("duration", "-duration=1m")
	if err != nil {
		t.Fatal(err)
	}
	if d != (60 * time.Second) {
		t.Fatal("invalid duration parsed")
	}
}
