package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/config"
)

// CORS -.
func CORS(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.CORS.AllowOrigins, ","), // Gin expects a slice
		AllowMethods:     strings.Split(cfg.CORS.AllowMethods, ","),
		AllowHeaders:     strings.Split(cfg.CORS.AllowHeaders, ","),
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: cfg.CORS.AllowCredentials,
		AllowOriginFunc:  func(origin string) bool { return true }, // Or implement custom logic if needed
		MaxAge:           12 * time.Hour,
	})
}
