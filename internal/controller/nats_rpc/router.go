package v1

import (
	v1 "github.com/minhhoccode111/realworld-fiber-clean/internal/controller/nats_rpc/v1"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/nats/nats_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation, l logger.Interface) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)

	{
		v1.NewTranslationRoutes(routes, t, l)
	}

	return routes
}
