package http

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"librenote/app/model"
	"librenote/app/response"
	"time"
)

type jwtCustomClaims struct {
	ID int32 `json:"id"`
	jwt.StandardClaims
}

//// generate token
//jwtCfg := config.Get().Jwt
//claims := &jwtCustomClaims{
//user.ID,
//jwt.StandardClaims{
//ExpiresAt: time.Now().Add(jwtCfg.ExpireTime * time.Second).Unix(),
//},
//}
//
//unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//token, err = unsignedToken.SignedString([]byte(jwtCfg.SecretKey))
//if err != nil {
//return "", err
//}
//return token, err

// UserHandler represent the http handler for user
type UserHandler struct {
	UUseCase model.UserUsecase
}

func NewUserHandler(e *echo.Echo, us model.UserUsecase) {
	handler := &UserHandler{
		UUseCase: us,
	}

	v1 := e.Group("/api/v1")
	v1.POST("/registration", handler.Registration)
}

func (u *UserHandler) Registration(c echo.Context) error {
	var regReq registrationReq
	err := c.Bind(&regReq)
	if err != nil {
		return c.JSON(response.RespondError(response.ErrUnprocessableEntity, err))
	}

	if ok, err := isRegistrationReqValid(&regReq); !ok {
		valErr, errors := formatValidationError(err)
		if valErr != nil {
			return c.JSON(response.RespondError(response.ErrBadRequest, valErr))
		}
		return c.JSON(response.RespondValidationError(response.ErrBadRequest, errors))
	}

	user := model.User{
		FullName:  regReq.FullName,
		Email:     regReq.Email,
		Hash:      regReq.Password,
		IsActive:  1,
		UpdatedAt: time.Now().UTC(),
	}
	ctx := c.Request().Context()
	err = u.UUseCase.Registration(ctx, &user)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondSuccess("registration successful", nil))

}
