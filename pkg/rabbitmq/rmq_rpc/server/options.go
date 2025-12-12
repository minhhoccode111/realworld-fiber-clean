package server

import "time"

// Option -.
// Option configures a RabbitMQ RPC server.
type Option func(*Server)

// Timeout -.
// Timeout sets the timeout for processing a RabbitMQ RPC message.
func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// ConnWaitTime -.
// ConnWaitTime sets the time to wait between connection attempts.
func ConnWaitTime(timeout time.Duration) Option {
	return func(s *Server) {
		s.conn.WaitTime = timeout
	}
}

// ConnAttempts -.
// ConnAttempts sets the number of attempts to connect to RabbitMQ.
func ConnAttempts(attempts int) Option {
	return func(s *Server) {
		s.conn.Attempts = attempts
	}
}
