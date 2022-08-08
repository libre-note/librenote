package http

import (
	"errors"
	"librenote/app/model"
	"librenote/app/response"
	"librenote/app/validation"
	"librenote/infrastructure/config"
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
	me.POST("", handler.UpdateSettings)
	me.DELETE("", handler.DeleteMe)
}

func (u *UserHandler) Registration(c echo.Context) error {
	if !config.Get().App.RegistrationOpen {
		return c.JSON(response.RespondError(response.ErrBadRequest, errors.New("registration closed")))
	}

	var regReq registrationReq

	err := c.Bind(&regReq)
	if err != nil {
		return c.JSON(response.RespondError(response.ErrUnprocessableEntity, err))
	}

	if ok, err := validation.Validate(&regReq); !ok {
		valErrors, valErr := validation.FormatErrors(err)
		if valErr != nil {
			return c.JSON(response.RespondError(response.ErrBadRequest, valErr))
		}

		return c.JSON(response.RespondValidationError(response.ErrBadRequest, valErrors))
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
		valErrors, valErr := validation.FormatErrors(err)
		if valErr != nil {
			return c.JSON(response.RespondError(response.ErrBadRequest, valErr))
		}

		return c.JSON(response.RespondValidationError(response.ErrBadRequest, valErrors))
	}

	ctx := c.Request().Context()

	token, err := u.UUseCase.Login(ctx, lReq.Email, lReq.Password)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondLoginSuccess(token))
}

func (u *UserHandler) Me(c echo.Context) error {
	ctx := c.Request().Context()

	details, err := u.UUseCase.GetUserDetails(ctx, getUserID(c))
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondSuccess("request success", details))
}

func getUserID(c echo.Context) int32 {
	token := c.Get("user").(*jwt.Token)
	return token.Claims.(*middlewares.JwtCustomClaims).UserID
}

func (u *UserHandler) UpdateSettings(c echo.Context) error {
	var usReq updateSettings

	err := c.Bind(&usReq)
	if err != nil {
		return c.JSON(response.RespondError(response.ErrUnprocessableEntity, err))
	}

	if ok, err := validation.Validate(&usReq); !ok {
		valErrors, valErr := validation.FormatErrors(err)
		if valErr != nil {
			return c.JSON(response.RespondError(response.ErrBadRequest, valErr))
		}

		return c.JSON(response.RespondValidationError(response.ErrBadRequest, valErrors))
	}

	// validate both password present
	if (len(usReq.OldPassword) > 0 && len(usReq.NewPassword) == 0) ||
		(len(usReq.NewPassword) > 0 && len(usReq.OldPassword) == 0) {
		return c.JSON(response.RespondError(
			response.ErrBadRequest,
			errors.New("both old_password & new_password field needed"),
		))
	}

	ctx := c.Request().Context()

	user, err := u.UUseCase.GetUser(ctx, getUserID(c))
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	passwordM := model.Password{
		OldPassword: usReq.OldPassword,
		NewPassword: usReq.NewPassword,
	}

	if len(usReq.OldPassword) > 0 && len(usReq.NewPassword) > 0 {
		passwordM.IsChanged = true
	}

	user.ListViewEnabled = *usReq.ListViewEnabled
	user.DarkModeEnabled = *usReq.DarkModeEnabled

	err = u.UUseCase.Update(ctx, user, passwordM)
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondSuccess("updated successfully", nil))
}

func (u *UserHandler) DeleteMe(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := u.UUseCase.GetUser(ctx, getUserID(c))
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	user.IsTrashed = 1

	err = u.UUseCase.Update(ctx, user, model.Password{})
	if err != nil {
		return c.JSON(response.RespondError(err))
	}

	return c.JSON(response.RespondEmpty())
}
