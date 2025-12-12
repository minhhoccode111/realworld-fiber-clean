package client

import "time"

// Option -.
// Option configures a NATS RPC client.
type Option func(*Client)

// Timeout -.
// Timeout sets the timeout for NATS RPC requests.
func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}
