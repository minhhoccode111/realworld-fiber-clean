package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

type ctxKey string

const (
	CtxUserIdKey ctxKey = "userId"
	CtxIsAuthKey ctxKey = "isAuth"
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(fiber.Map{"error": msg})
}

// AuthMiddleware -.
func AuthMiddleware(l logger.Interface, jwtSecret string, isOptional bool) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals(CtxIsAuthKey, false)
		c.Locals(CtxUserIdKey, "")

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			if isOptional {
				return c.Next()
			}
			return errorResponse(c, http.StatusUnauthorized, "missing authorization header")
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 {
			if isOptional {
				return c.Next()
			}
			return errorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
		}
		if !strings.EqualFold(parts[0], "Token") {
			if isOptional {
				return c.Next()
			}
			return errorResponse(
				c,
				http.StatusUnauthorized,
				"authorization header must use 'Token' scheme",
			)
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			if isOptional {
				return c.Next()
			}
			l.Error(err, "http - middleware - AuthMiddleware - jwt.Parse")
			return errorResponse(c, http.StatusUnauthorized, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			if isOptional {
				return c.Next()
			}
			return errorResponse(c, http.StatusUnauthorized, "invalid token claims")
		}

		userId, ok := claims["sub"].(string)
		if !ok || userId == "" {
			if isOptional {
				return c.Next()
			}
			return errorResponse(c, http.StatusUnauthorized, "missing user id in token")
		}

		c.Locals(CtxIsAuthKey, true)
		c.Locals(CtxUserIdKey, userId)
		return c.Next()
	}
}
