package validation

import (
	"regexp"
	"strings"
	"web_socket/internal/utils"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidation(v *validator.Validate) {
	var blockedDomains = map[string]bool{
		"hotmail.com": true,
		"abc.com":     true,
	}

	v.RegisterValidation("email_advanced", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return false
		}

		domain := utils.NormalizeString(parts[1])
		return !blockedDomains[domain]
	})

	v.RegisterValidation("password_string", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}

		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>?/\\|]`).MatchString(password)

		return hasDigit && hasLower && hasUpper && hasSpecial
	})

	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})
}
