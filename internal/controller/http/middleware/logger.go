package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

func buildRequestMessage(ctx *gin.Context) string {
	var result strings.Builder

	result.WriteString(ctx.ClientIP())
	result.WriteString(" - ")
	result.WriteString(ctx.Request.Method)
	result.WriteString(" ")
	result.WriteString(ctx.Request.URL.Path)
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(ctx.Writer.Status()))
	// Gin does not directly expose response body length in ctx.Writer
	// result.WriteString(" ")
	// result.WriteString(strconv.Itoa(len(ctx.Response().Body()))) // This was Fiber specific

	return result.String()
}

func Logger(l logger.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		l.Info(buildRequestMessage(ctx))
	}
}
