package core

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func formatValidationError(err error) map[string]string {
	out := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()

			switch fe.Tag() {
			case "required":
				out[field] = "This field is required"
			case "min":
				out[field] = fmt.Sprintf("Must be at least %s characters long", fe.Param())
			case "max":
				out[field] = fmt.Sprintf("Cannot exceed %s characters", fe.Param())
			case "e164":
				out[field] = "Must be a valid international phone number"
			case "email":
				out[field] = "Must be a valid email address"
			case "latitude":
				out[field] = "Must be between -90 and 90."
			case "longitude":
				out[field] = "Must be between -180 and 180."
			case "oneof":
				out[field] = fmt.Sprintf("Must be one of: %s", fe.Param())
			case "gt":
				out[field] = fmt.Sprintf("Must be greater than %s", fe.Param())
			case "gte":
				out[field] = fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
			case "lt":
				out[field] = fmt.Sprintf("Must be less than %s", fe.Param())
			case "lte":
				out[field] = fmt.Sprintf("Must be less than or equal to %s", fe.Param())
			default:
				out[field] = fmt.Sprintf("Failed validation on rule '%s'", fe.Tag())
			}
		}
		return out
	}

	out["body"] = "Invalid JSON syntax or invalid request body"
	return out
}
