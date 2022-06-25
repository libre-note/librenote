package http

import (
	"errors"
	"fmt"
	validator "github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type registrationReq struct {
	FullName string `json:"full_name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func isRegistrationReqValid(r *registrationReq) (bool, error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	err := validate.Struct(r)
	if err != nil {
		return false, err
	}
	return true, nil
}

func formatValidationError(err error) (error, map[string]interface{}) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]interface{}, len(ve))
		for _, fe := range ve {
			fieldName := fe.Field()
			out[fieldName] = msgForField(fe)
		}
		return nil, out
	}
	return err, nil
}

func msgForField(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		m := fmt.Sprintf("Must be at least %v", fe.Param())
		if fe.Type().String() == "string" {
			return m + " characters"
		}
		return m
	case "max":
		m := fmt.Sprintf("Not more than %v", fe.Param())
		if fe.Type().String() == "string" {
			return m + " characters"
		}
		return m
	}

	return "unknown error"
}
