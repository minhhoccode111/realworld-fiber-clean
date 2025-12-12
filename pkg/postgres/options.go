package postgres

import "time"

// Option -.
// Option configures a Postgres connection.
type Option func(*Postgres)

// MaxPoolSize -.
// MaxPoolSize sets the maximum number of connections in the pool.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
// ConnAttempts sets the number of attempts to connect to the Postgres database.
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
// ConnTimeout sets the timeout for each connection attempt to the Postgres database.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}
