package server

import "time"

// Option -.
// Option configures a NATS RPC server.
type Option func(*Server)

// Timeout -.
// Timeout sets the timeout for NATS RPC message processing.
func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}
