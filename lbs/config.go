package lbs

import (
	"errors"
	"time"
)

const defaultTimeout = 10 * time.Second

// Config Tencent Location Service API configuration.
type Config struct {
	// AppKey is the Tencent LBS API key.
	AppKey string
	// Timeout is the HTTP client timeout; defaults to 10s when zero.
	Timeout time.Duration
}

// Validate checks required fields.
func (c Config) Validate() error {
	if c.AppKey == "" {
		return errors.New("lbs: app key is required")
	}
	return nil
}

func (c Config) timeout() time.Duration {
	if c.Timeout > 0 {
		return c.Timeout
	}
	return defaultTimeout
}
