package main

import (
	"errors"
	"strings"
	"time"
)

func parseSwitch(name string, arg string) (time.Duration, error) {
	if arg[0] != '-' {
		return 0, errors.New("invalid switch")
	}
	pre := "-"+name+"="
	if !strings.HasPrefix(arg, pre) {
		return 0, errors.New("switch did not match")
	}
	s := arg[len(pre):]
	return time.ParseDuration(s)
}
