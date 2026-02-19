package response

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	errs := make(map[string]string)

	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, e := range ve {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errs[field] = field + " is required"
			case "min":
				errs[field] = field + " must be at least " + e.Param() + " characters"
			case "email":
				errs[field] = field + " must be a valid email address"
			case "e164":
				errs[field] = field + " must be a valid E.164 phone number"
			default:
				errs[field] = "Invalid value for " + field
			}
		}
	}

	return errs
}
