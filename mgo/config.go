package mgo

import (
	"time"
)

// GlobalConfig encapsulates the global configuration to be used by mgo
type GlobalConfig struct {
	// DialMaxRetries defines the maximum amount of retries to attempt when dialing to a
	// connection to a mongodb instance
	DialMaxRetries int `json:"dial_max_retries" yaml:"dial_max_retries"`

	// DialRetryTimeout defines the timeout in milliseconds between retries when dialing
	// for a connection to a mongodb instance.
	DialRetryTimeout time.Duration `json:"dial_retry_timeout" yaml:"dial_max_retries"`
}

var instance *GlobalConfig

// Config retrieves the global configuration used by the mgo package
func Config() *GlobalConfig {
	if instance == nil {
		instance = &GlobalConfig{
			DialMaxRetries:   3,
			DialRetryTimeout: 10000 * time.Millisecond,
		}
	}
	return instance
}
