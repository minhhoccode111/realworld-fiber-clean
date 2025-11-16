package validatorx

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	_ = v.RegisterValidation("no_dups_str", func(fl validator.FieldLevel) bool {
		slices, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}

		seen := make(map[string]struct{})
		for _, t := range slices {
			t = strings.TrimSpace(t)
			if _, exists := seen[t]; exists {
				return false
			}
			seen[t] = struct{}{}
		}

		return true
	})

	_ = v.RegisterValidation("tag", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z0-9_ -]+$`).MatchString(fl.Field().String())
	})

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
