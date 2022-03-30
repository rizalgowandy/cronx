// TODO: replace me
package library

import (
	"github.com/rizalgowandy/library-template-go/pkg/api"
)

// TODO: replace me
// NewClient creates a client to interact with XYZ API.
func NewClient(cfg api.Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Client{}, nil
}

// TODO: replace me
// Client is the main client to interact with XYZ API.
type Client struct {
}
