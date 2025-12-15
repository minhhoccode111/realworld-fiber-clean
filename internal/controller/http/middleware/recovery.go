package middleware

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

func buildPanicMessage(ctx *fiber.Ctx, err any) string {
	var result strings.Builder

	result.WriteString(ctx.IP())
	result.WriteString(" - ")
	result.WriteString(ctx.Method())
	result.WriteString(" ")
	result.WriteString(ctx.OriginalURL())
	result.WriteString(" PANIC DETECTED: ")
	result.WriteString(fmt.Sprintf("%v\n%s\n", err, debug.Stack()))

	return result.String()
}

// logPanic returns a handler that logs panic details using the provided logger.
func logPanic(l logger.Interface) func(c *fiber.Ctx, err any) {
	return func(ctx *fiber.Ctx, err any) {
		l.Error(buildPanicMessage(ctx, err))
	}
}

// Recovery installs a panic recovery middleware that logs stack traces.
func Recovery(l logger.Interface) func(c *fiber.Ctx) error {
	return fiberRecover.New(fiberRecover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: logPanic(l),
	})
}
