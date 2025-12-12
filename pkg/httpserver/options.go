package httpserver

import (
	"net"
	"time"
)

// Option -.
// Option configures a Server.
type Option func(*Server)

// Port -.
// Port sets the address port for the HTTP server.
func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

// Prefork -.
// Prefork enables or disables Fiber's prefork feature.
func Prefork(prefork bool) Option {
	return func(s *Server) {
		s.prefork = prefork
	}
}

// ReadTimeout -.
// ReadTimeout sets the maximum duration for reading the entire
// request, including the body.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WriteTimeout -.
// WriteTimeout sets the maximum duration before timing out
// writes of the response.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}

// ShutdownTimeout -.
// ShutdownTimeout sets the maximum duration for the server
// to gracefully shut down.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
