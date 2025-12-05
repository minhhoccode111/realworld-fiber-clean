package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

func buildPanicMessage(ctx *gin.Context, err any) string {
	var result strings.Builder

	result.WriteString(ctx.ClientIP())
	result.WriteString(" - ")
	result.WriteString(ctx.Request.Method)
	result.WriteString(" ")
	result.WriteString(ctx.Request.URL.Path)
	result.WriteString(" PANIC DETECTED: ")
	result.WriteString(fmt.Sprintf("%v\n%s\n", err, debug.Stack()))

	return result.String()
}

func logPanic(l logger.Interface) gin.RecoveryFunc {
	return func(ctx *gin.Context, err any) {
		l.Error(buildPanicMessage(ctx, err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

func Recovery(l logger.Interface) gin.HandlerFunc {
	return gin.CustomRecovery(logPanic(l))
}
