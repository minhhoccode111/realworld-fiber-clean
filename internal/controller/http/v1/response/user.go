package response

import (
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// UserAuth -.
type UserAuth struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func NewUserAuth(u *entity.User, token string) UserAuth {
	return UserAuth{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Token:    token,
	}
}

// UserAuthResponse -.
type UserAuthResponse struct {
	User UserAuth `json:"user"`
}
