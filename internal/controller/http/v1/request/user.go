package request

import "strings"

// UserRegister -.
type UserRegister struct {
	Email    string `json:"email"    validate:"required,email"                 example:"minhhoccode111@gmail.com"`
	Username string `json:"username" validate:"required,min=2,max=50,username" example:"minhhoccode111"`
	Password string `json:"password" validate:"required,min=8,max=50,password" example:"P@ssw0rd"`
}

func (ur *UserRegister) Trim() {
	ur.Email = strings.TrimSpace(ur.Email)
	ur.Username = strings.TrimSpace(ur.Username)
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
