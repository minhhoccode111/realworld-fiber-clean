package request

import "strings"

// UserRegister -.
type UserRegister struct {
	Email    string `json:"email"    validate:"required,email"                 example:"minhhoccode111@gmail.com"`
	Username string `json:"username" validate:"required,min=2,max=50,username" example:"minhhoccode111"`
	Password string `json:"password" validate:"required,min=8,max=50,password" example:"P@ssw0rd"`
}

func (u *UserRegister) Trim() {
	u.Email = strings.TrimSpace(u.Email)
	u.Username = strings.TrimSpace(u.Username)
	// ur.Password = strings.TrimSpace(ur.Password) // WARN: don't trim password
}

type UserRegisterRequest struct {
	User UserRegister `json:"user"`
}

// UserLogin -.
type UserLogin struct {
	Email    string `json:"email"    validate:"required" example:"minhhoccode111@gmail.com"`
	Password string `json:"password" validate:"required" example:"P@ssw0rd"`
}

func (ur *UserLogin) Trim() {
	ur.Email = strings.TrimSpace(ur.Email)
}

type UserLoginRequest struct {
	User UserLogin `json:"user"`
}

// UserUpdate -.
type UserUpdate struct {
	Email    string `json:"email"    validate:"omitempty,min=5,max=320,email"   example:"minhhoccode111@gmail.com"`
	Username string `json:"username" validate:"omitempty,min=2,max=50,username" example:"minhhoccode111"`
	Password string `json:"password" validate:"omitempty,min=8,max=50,password" example:"P@ssw0rd"`
	Bio      string `json:"bio"      validate:"max=255"                         example:"Trust the process"`
	Image    string `json:"image"    validate:"max=2048"                        example:"https://www.w3schools.com/howto/img_avatar.png"`
}

func (u *UserUpdate) Trim() {
	u.Email = strings.TrimSpace(u.Email)
	u.Username = strings.TrimSpace(u.Username)
	u.Bio = strings.TrimSpace(u.Bio)
	u.Image = strings.TrimSpace(u.Image)
}

type UserUpdateRequest struct {
	User UserUpdate `json:"user"`
}
