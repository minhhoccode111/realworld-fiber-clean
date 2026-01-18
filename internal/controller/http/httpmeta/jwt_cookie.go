package httpmeta

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// NewJWTInCookie constructs a JWT-bearing cookie with the provided TTL.
func NewJWTInCookie(token string, duration time.Duration, secure bool) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     CookieJWTName,
		Value:    token,
		Expires:  time.Now().Add(duration),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
		Path:     "/",
	}
}
