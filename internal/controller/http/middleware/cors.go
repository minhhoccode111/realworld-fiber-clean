package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
)

// CORS -.
func CORS(cfg *config.Config) func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.CORS_ALLOW_ORIGINS,
		AllowHeaders:     cfg.CORS.CORS_ALLOW_HEADERS,
		AllowCredentials: cfg.CORS.CORS_ALLOW_CREDENTIALS,
		AllowMethods:     cfg.CORS.CORS_ALLOW_METHODS,
	})
}
