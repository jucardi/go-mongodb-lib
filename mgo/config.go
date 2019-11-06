package mgo

import (
	"time"
)

var (
	// DialMaxRetries defines the maximum amount of retries to attempt when dialing to a
	// connection to a mongodb instance
	DialMaxRetries = 3

	// DialRetrySleep defines the sleep time between retries when dialing for a connection to a mongodb instance.
	DialRetrySleep = 10 * time.Second

	// DialTimeout indicates the max time to wait before aborting a dialing attempt.
	DialTimeout = 10 * time.Second
)
