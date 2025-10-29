// Package v1 implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
	_ "github.com/minhhoccode111/realworld-fiber-clean/docs" // Swagger docs.
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	v1 "github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Realworld Fiber Clean API
// @description Realworld API using Golang + Fiber + Clean Architecture
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Token" followed by a space and JWT token.
func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	l logger.Interface,

	t usecase.Translation,
	tc usecase.TranslationClone,
	u usecase.User,
	a usecase.Article,
	f usecase.Favorite,
	c usecase.Comment,
	tag usecase.Tag,
) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New(cfg.App.Name)
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("/api/v1")
	{
		v1.NewV1Routes(apiV1Group, cfg, l, t, tc, u, a, f, c, tag)
	}
}
