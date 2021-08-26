package main

import (
	"github.com/juju/ratelimit"
	"time"
)

var endpoints map[string]*ratelimit.Bucket

func shouldThrottle(host string) bool {
	bucket, ok := endpoints[host]
	if !ok {
		bucket = ratelimit.NewBucket(time.Minute, 300)
	}
	ok = bucket.WaitMaxDuration(1, 3 * time.Second)
	endpoints[host] = bucket
	return !ok // If we're ok, we shouldn't throttle
}

func init() {
	endpoints = make(map[string]*ratelimit.Bucket)
}
