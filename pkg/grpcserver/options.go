package grpcserver

import (
	"net"
)

// Option -.
// Option configures a Server.
type Option func(*Server)

// Port -.
// Port sets the address port for the gRPC server.
func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}
