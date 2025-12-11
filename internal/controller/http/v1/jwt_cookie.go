package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/common"
)

func NewJWTCookie(token string, duration time.Duration) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     common.CookieJWTName,
		Value:    token,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	}
}
