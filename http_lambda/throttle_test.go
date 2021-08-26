package main

import (
	"testing"
	"time"
)

func TestThrottling(t *testing.T) {
	host := "104.62.139.94"
	if shouldThrottle(host) {
		t.Error("should not have throttled a single request")
	}

	didWeThrottle := false
	for i := 0; i < 500; i++ {
		if shouldThrottle(host) {
			didWeThrottle = true
		}
		time.Sleep(2 * time.Millisecond)
	}
	if !didWeThrottle {
		t.Error("should have throttled with that many requests")
	}
}
