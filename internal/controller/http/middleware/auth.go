package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/common"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(fiber.Map{"error": msg})
}

// AuthMiddleware -.
func AuthMiddleware(l logger.Interface, jwtSecret string, isOptional bool) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		n := 2
		s := fmt.Sprintf(
			"Authorization header must be formatted: [%s <token>]",
			common.AuthorizationScheme,
		)
		esc := func(msg string) error {
			if isOptional {
				return c.Next()
			}

			return errorResponse(c, http.StatusUnauthorized, msg)
		}

		c.Locals(common.CtxIsAuthKey, false)
		c.Locals(common.CtxUserIDKey, "")
		c.Locals(common.CtxUserRoleKey, "")

		var tokenStr string

		authHeader := c.Get("Authorization")
		if authHeader != "" {
			// use jwt-in-header
			parts := strings.Fields(authHeader)
			if len(parts) < n {
				return esc(s)
			}

			if !strings.EqualFold(parts[0], common.AuthorizationScheme) {
				return esc(s)
			}

			tokenStr = parts[1]
		} else {
			// use jwt-in-cookie
			tokenStr = c.Cookies(common.CookieJWTName)
			if tokenStr == "" {
				return esc("No token provided")
			}
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%w: %v", entity.ErrUnexpectedSigningMethod, t.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})
		if err != nil {
			l.Error(err, "http - middleware - AuthMiddleware - jwt.Parse")

			return esc("Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return esc("Invalid token claims")
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			return esc("Missing userID in token")
		}

		roleStr, ok := claims["role"].(string)

		userRole := entity.Role(roleStr)
		if !userRole.IsValid() {
			userRole = entity.UserRole
		}

		c.Locals(common.CtxIsAuthKey, true)
		c.Locals(common.CtxUserIDKey, userID)
		c.Locals(common.CtxUserRoleKey, userRole)

		return c.Next()
	}
}
