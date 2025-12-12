package httpmeta

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewJWTInCookie(token string, duration time.Duration) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     CookieJWTName,
		Value:    token,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	}
}
