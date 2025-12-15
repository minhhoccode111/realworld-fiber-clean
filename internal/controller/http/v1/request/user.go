package request

import (
	"strings"

	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// UserRegister describes the payload for user registration.
type UserRegister struct {
	Email    string `json:"email"    validate:"required,email"                 example:"minhhoccode111@gmail.com"`
	Username string `json:"username" validate:"required,min=2,max=50,username" example:"minhhoccode111"`
	Password string `json:"password" validate:"required,min=8,max=50,password" example:"P@ssw0rd"`
}

func (u *UserRegister) Trim() {
	u.Email = strings.TrimSpace(u.Email)
	u.Username = strings.TrimSpace(u.Username)
}

type UserRegisterRequest struct {
	User UserRegister `json:"user"`
}

// UserLogin describes the payload for user login.
type UserLogin struct {
	Email    string `json:"email"    validate:"required" example:"minhhoccode111@gmail.com"`
	Password string `json:"password" validate:"required" example:"P@ssw0rd"`
}

func (u *UserLogin) Trim() {
	u.Email = strings.TrimSpace(u.Email)
}

type UserLoginRequest struct {
	User UserLogin `json:"user"`
}

// UserUpdate holds optional fields for updating a user.
type UserUpdate struct {
	Email    string `json:"email"    validate:"omitempty,min=5,max=320,email"   example:"minhhoccode111@gmail.com"`
	Username string `json:"username" validate:"omitempty,min=2,max=50,username" example:"minhhoccode111"`
	Password string `json:"password" validate:"omitempty,min=8,max=50,password" example:"P@ssw0rd"`
	Bio      string `json:"bio"      validate:"max=255"                         example:"Trust the process"`
	Image    string `json:"image"    validate:"max=2048"                        example:"https://www.w3schools.com/howto/img_avatar.png"`
}

// NewUser converts an update payload into a full user entity.
func (uu *UserUpdate) NewUser(userID string) *entity.User {
	return &entity.User{
		ID:       userID,
		Email:    uu.Email,
		Username: uu.Username,
		Image:    uu.Image,
		Bio:      uu.Bio,
		Password: uu.Password,
	}
}

func (uu *UserUpdate) Trim() {
	uu.Email = strings.TrimSpace(uu.Email)
	uu.Username = strings.TrimSpace(uu.Username)
	uu.Bio = strings.TrimSpace(uu.Bio)
	uu.Image = strings.TrimSpace(uu.Image)
}

func (uu *UserUpdate) IsAllEmpty() bool {
	b := uu.Bio == ""
	e := uu.Email == ""
	i := uu.Image == ""
	u := uu.Username == ""
	p := uu.Password == ""

	return e && u && b && i && p
}

type UserUpdateRequest struct {
	User UserUpdate `json:"user"`
}
