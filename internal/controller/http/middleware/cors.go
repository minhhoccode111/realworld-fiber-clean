package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
)

// CORS configures CORS middleware from application config.
func CORS(cfg *config.Config) func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		AllowMethods:     cfg.CORS.AllowMethods,
	})
}
