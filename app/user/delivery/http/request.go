package http

import (
	validator "github.com/go-playground/validator/v10"
)

type registrationReq struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func isRegistrationReqValid(r *registrationReq) (bool, error) {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return false, err
	}
	return true, nil
}
