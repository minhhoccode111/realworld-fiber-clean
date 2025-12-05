// Package v1 implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
	_ "github.com/minhhoccode111/realworld-fiber-clean/docs" // Swagger docs.
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	v1 "github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/usecase"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
	"github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zsais/go-gin-prometheus"
)

// NewRouter -.
// Swagger spec:
// @title       Realworld Fiber Clean API
// @description Realworld API using Golang + Gin + Clean Architecture
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Token" followed by a space and JWT token.
func NewRouter(
	router *gin.Engine,
	cfg *config.Config,
	l logger.Interface,

	t usecase.Translation,
	u usecase.User,
	a usecase.Article,
	f usecase.Favorite,
	c usecase.Comment,
	p usecase.Profile,
	tag usecase.Tag,
) {
	// Options
	router.Use(middleware.Logger(l))
	router.Use(middleware.Recovery(l))
	router.Use(middleware.CORS(cfg))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		p := ginprometheus.NewPrometheus("gin")
		p.Use(router)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// K8s probe
	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	apiV1Group := router.Group("/api/v1")
	{
		v1.NewV1Routes(apiV1Group, cfg, l, t, u, a, f, c, p, tag)
	}
}
