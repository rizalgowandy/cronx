package api

import (
	"errors"
	"time"
)

// Config is the necessary configuration to call API.
type Config struct {
	// Key is the authentication API key.
	// Most requests to the XYZ API must be authenticated with an API key.
	// You can create an API key in your Settings page after creating a XYZ account.
	// Reference: TODO: replace
	Key string
	// Timeout describes total waiting time before a request is treated as timeout.
	// Optional.
	// Default: 1 min.
	Timeout time.Duration
	// RetryCount describes total number of retry in case error occurred.
	// Optional.
	// Default: 0 = disable retry mechanism.
	RetryCount int
	// RetryMaxWaitTime describes total waiting time between each retry.
	// Optional.
	// Default: 2 second.
	RetryMaxWaitTime time.Duration
	// Debug describes the client to enter debug mode.
	// Debug mode will dump the request and response on each API call.
	// Be warn, credentials data will be dumped too.
	// Ensure you're only this mode on safe environment like local.
	// Optional.
	// Default: false.
	Debug bool
	// HostURL describes the host url target.
	// HostURL can be filled with your fake server host url for testing purpose.
	// Optional.
	// Default: TODO: replace
	HostURL string
}

// Validate validates configuration correctness and
// fill fields with default configuration if left empty.
func (c *Config) Validate() error {
	if c.Key == "" {
		return errors.New("config: invalid key")
	}
	if c.Timeout <= 0 {
		c.Timeout = time.Minute
	}
	if c.RetryCount < 0 {
		c.RetryCount = 3
	}
	if c.RetryMaxWaitTime <= 0 {
		c.RetryMaxWaitTime = 2 * time.Second
	}
	if c.HostURL == "" {
		c.HostURL = "" // TODO: replace
	}
	return nil
}
