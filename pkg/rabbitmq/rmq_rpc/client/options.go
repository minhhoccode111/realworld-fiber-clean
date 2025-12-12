package client

import "time"

// Option -.
// Option configures a RabbitMQ RPC client.
type Option func(*Client)

// Timeout -.
// Timeout sets the timeout for RabbitMQ RPC requests.
func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// ConnWaitTime -.
// ConnWaitTime sets the time to wait between connection attempts.
func ConnWaitTime(timeout time.Duration) Option {
	return func(c *Client) {
		c.conn.WaitTime = timeout
	}
}

// ConnAttempts -.
// ConnAttempts sets the number of attempts to connect to RabbitMQ.
func ConnAttempts(attempts int) Option {
	return func(c *Client) {
		c.conn.Attempts = attempts
	}
}
