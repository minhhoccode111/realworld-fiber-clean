package rmqrpc

import "errors"

var (
	// ErrTimeout -.
	// ErrTimeout indicates that a RabbitMQ RPC request timed out.
	ErrTimeout = errors.New("timeout")
	// ErrInternalServer -.
	// ErrInternalServer indicates an internal server error occurred in the RabbitMQ RPC server.
	ErrInternalServer = errors.New("internal server error")
	// ErrBadHandler -.
	// ErrBadHandler indicates that the requested handler is not registered on the RabbitMQ RPC server.
	ErrBadHandler = errors.New("unregistered handler")
)

// Success -.
// Success is the status string for a successful RabbitMQ RPC call.
const Success = "success"
