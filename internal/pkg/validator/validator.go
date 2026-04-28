package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return Validate(dst)
}

func Validate(dst any) error {
	if err := validate.Struct(dst); err != nil {
		var errs validator.ValidationErrors
		if ok := errors.As(err, &errs); ok {
			msgs := make([]string, 0, len(errs))
			for _, e := range errs {
				msgs = append(msgs, fieldError(e))
			}
			return fmt.Errorf("%s", strings.Join(msgs, "; "))
		}
		return err
	}
	return nil
}

func fieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", e.Field())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
