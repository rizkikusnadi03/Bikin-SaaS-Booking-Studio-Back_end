package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	errs := make(map[string]string)
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				errs[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errs[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				errs[field] = fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
			case "max":
				errs[field] = fmt.Sprintf("%s must not exceed %s characters", field, fe.Param())
			case "uuid":
				errs[field] = fmt.Sprintf("%s must be a valid UUID", field)
			default:
				errs[field] = fmt.Sprintf("%s failed validation on '%s'", field, fe.Tag())
			}
		}
		return errs
	}

	errs["general"] = err.Error()
	return errs
}
