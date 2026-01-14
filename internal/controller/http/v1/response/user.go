package response

import (
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// UserAuth represents user data returned after authentication.
type UserAuth struct {
	Email    string      `json:"email"`
	Username string      `json:"username"`
	Bio      string      `json:"bio"`
	Image    string      `json:"image"`
	Token    string      `json:"token"`
	Role     entity.Role `json:"role"`
}

// NewUserAuth builds a UserAuth response from a user entity and JWT token.
func NewUserAuth(u *entity.User, token string) UserAuth {
	return UserAuth{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Role:     u.Role,
		Token:    token,
	}
}

// UserAuthResponse wraps authenticated user information.
type UserAuthResponse struct {
	User UserAuth `json:"user"`
}
