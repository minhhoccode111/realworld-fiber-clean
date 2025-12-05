package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
)

type ctxKey string

const (
	CtxUserIDKey   ctxKey = "userID"
	CtxUserRoleKey ctxKey = "userRole"
	CtxIsAuthKey   ctxKey = "isAuth"
)

func errorResponse(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(code, gin.H{"error": msg})
}

// AuthMiddleware -.
func AuthMiddleware(l logger.Interface, jwtSecret string, isOptional bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(string(CtxIsAuthKey), false)
		c.Set(string(CtxUserIDKey), "")
		c.Set(string(CtxUserRoleKey), "")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if isOptional {
				c.Next()
				return
			}

			errorResponse(c, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Fields(authHeader)

		const lenParts = 2
		if len(parts) != lenParts {
			if isOptional {
				c.Next()
				return
			}

			errorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		if !strings.EqualFold(parts[0], "Token") {
			if isOptional {
				c.Next()
				return
			}

			errorResponse(
				c,
				http.StatusUnauthorized,
				"authorization header must use 'Token' scheme",
			)
			return
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
				c.Next()
				return
			}

			l.Error(err, "http - middleware - AuthMiddleware - jwt.Parse")

			errorResponse(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			if isOptional {
				c.Next()
				return
			}

			errorResponse(c, http.StatusUnauthorized, "invalid token claims")
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			if isOptional {
				c.Next()
				return
			}

			errorResponse(c, http.StatusUnauthorized, "missing user id in token")
			return
		}

		roleStr, ok := claims["role"].(string)
		userRole := entity.Role(roleStr)
		if !userRole.IsValid() {
			userRole = entity.UserRole
		}

		c.Set(string(CtxIsAuthKey), true)
		c.Set(string(CtxUserIDKey), userID)
		c.Set(string(CtxUserRoleKey), userRole)

		c.Next()
	}
}
