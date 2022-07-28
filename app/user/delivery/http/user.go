package http

import (
	"librenote/app/model"
	"librenote/app/response"
	"librenote/app/validation"
	"librenote/infrastructure/middlewares"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

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
	v1.POST("/login", handler.Login)

	me := e.Group("/api/v1/me")
	_ = middlewares.AttachJwtToGroup(me)
	me.GET("", handler.Me)
}

func (u *UserHandler) Registration(c echo.Context) error {
	var regReq registrationReq
	err := c.Bind(&regReq)
	if err != nil {
		return c.JSON(response.RespondError(response.ErrUnprocessableEntity, err))
	}

	if ok, err := validation.Validate(&regReq); !ok {
		valErr, errors := validation.FormatErrors(err)
		if valErr != nil {
			return c.JSON(response.RespondError(response.ErrBadRequest, valErr))
		}
		return c.JSON(response.RespondValidationError(response.ErrBadRequest, errors))
	}

	nowTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	user := model.User{
		FullName:  regReq.FullName,
		Email:     regReq.Email,
		Hash:      regReq.Password,
		IsActive:  1,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}
	ctx := c.Request().Context()
	err = u.UUseCase.Registration(ctx, &user)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondSuccess("registration successful", nil))
}

func (u *UserHandler) Login(c echo.Context) error {
	var lReq loginReq
	err := c.Bind(&lReq)
	if err != nil {
		return c.JSON(response.RespondError(response.ErrUnprocessableEntity, err))
	}

	if ok, err := validation.Validate(&lReq); !ok {
		valErr, errors := validation.FormatErrors(err)
		if valErr != nil {
			return c.JSON(response.RespondError(response.ErrBadRequest, valErr))
		}
		return c.JSON(response.RespondValidationError(response.ErrBadRequest, errors))
	}

	ctx := c.Request().Context()
	token, err := u.UUseCase.Login(ctx, lReq.Email, lReq.Password)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondLoginSuccess(token))
}

func (u *UserHandler) Me(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	userID := token.Claims.(*middlewares.JwtCustomClaims).UserID

	ctx := c.Request().Context()
	details, err := u.UUseCase.GetUserDetails(ctx, userID)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondSuccess("request success", details))
}
