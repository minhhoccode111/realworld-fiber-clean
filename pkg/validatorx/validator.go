package validatorx

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(fl.Field().String())
	})

	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		p := fl.Field().String()
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(p)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(p)
		hasDigit := regexp.MustCompile(`\d`).MatchString(p)
		hasSpecial := regexp.MustCompile(`[!@#~$%^&*()+|_{}<>?,./-]`).MatchString(p)
		return hasUpper && hasLower && hasDigit && hasSpecial
	})

	return v
}
