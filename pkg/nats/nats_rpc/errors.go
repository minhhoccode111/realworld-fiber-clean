package natsrpc

import "errors"

var (
	// ErrTimeout -.
	// ErrTimeout indicates that a NATS RPC request timed out.
	ErrTimeout = errors.New("timeout")
	// ErrInternalServer -.
	// ErrInternalServer indicates an internal server error occurred in the NATS RPC server.
	ErrInternalServer = errors.New("internal server error")
	// ErrBadHandler -.
	// ErrBadHandler indicates that the requested handler is not registered on the NATS RPC server.
	ErrBadHandler = errors.New("unregistered handler")
)

// Success -.
// Success is the status string for a successful NATS RPC call.
const Success = "success"
