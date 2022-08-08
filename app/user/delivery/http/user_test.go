package http_test

import (
	"encoding/json"
	"io"
	"librenote/app/model"
	"librenote/app/model/mocks"
	"librenote/app/response"
	userHttp "librenote/app/user/delivery/http"
	"librenote/infrastructure/config"
	"librenote/infrastructure/middlewares"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var BaseURLV1 = "/api/v1"

type registrationReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateSettings struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ListViewEnabled int8   `json:"list_view_enabled"`
	DarkModeEnabled int8   `json:"dark_mode_enabled"`
}

func TestRegistration(t *testing.T) {
	config.SetRegistrationOn()

	mockUsecase := new(mocks.UserUsecase)
	mockUsecase.On("Registration", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

	regReq := registrationReq{
		FullName: "Mr. Test",
		Email:    "mrtest@example.com",
		Password: "12345678",
	}

	endPoint := BaseURLV1 + "/registration"

	t.Run("success", func(t *testing.T) {
		tempReq := regReq
		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		c, rec := buildEchoPostRequest(t, endPoint, strings.NewReader(string(j)))

		handler := userHttp.UserHandler{
			UUseCase: mockUsecase,
		}
		err = handler.Registration(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("short password", func(t *testing.T) {
		tempReq := regReq
		tempReq.Password = "1234567"

		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		c, rec := buildEchoPostRequest(t, endPoint, strings.NewReader(string(j)))

		handler := userHttp.UserHandler{
			UUseCase: mockUsecase,
		}
		err = handler.Registration(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockUsecase := new(mocks.UserUsecase)
	mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("", nil)

	lReq := loginReq{
		Email:    "mrtest@example.com",
		Password: "12345678",
	}

	endPoint := BaseURLV1 + "/login"

	t.Run("success", func(t *testing.T) {
		tempReq := lReq
		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		c, rec := buildEchoPostRequest(t, endPoint, strings.NewReader(string(j)))

		handler := userHttp.UserHandler{
			UUseCase: mockUsecase,
		}
		err = handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("wrong email", func(t *testing.T) {
		tempReq := lReq
		tempReq.Email = "wrong-email.com"

		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		c, rec := buildEchoPostRequest(t, endPoint, strings.NewReader(string(j)))

		handler := userHttp.UserHandler{
			UUseCase: mockUsecase,
		}
		err = handler.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func buildEchoPostRequest(t *testing.T, path string, payload io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req, err := http.NewRequest(echo.POST, path, payload)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)

	return c, rec
}

func buildEchoAuthorizedRequest(t *testing.T, method, path, token string, payload io.Reader) (
	echo.Context, *httptest.ResponseRecorder) {
	var req *http.Request

	var err error

	if payload != nil {
		req, err = http.NewRequest(method, path, payload)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req, err = http.NewRequest(method, path, nil)
	}

	assert.NoError(t, err)

	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)

	res := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(req, res)

	return ctx, res
}

// nolint:unparam
func getToken(userID int32) string {
	jwtCfg := config.Get().Jwt
	claims := &middlewares.JwtCustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtCfg.ExpireTime).Unix(),
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := unsignedToken.SignedString([]byte(jwtCfg.SecretKey))

	return token
}

func attachJWTMiddleware(hfc echo.HandlerFunc) echo.HandlerFunc {
	mhfc := middleware.JWTWithConfig(
		middleware.JWTConfig{
			Claims:     &middlewares.JwtCustomClaims{},
			SigningKey: []byte(config.Get().Jwt.SecretKey),
		})(hfc)

	return mhfc
}

func TestMe(t *testing.T) {
	endPoint := BaseURLV1 + "/me"
	mockUser := model.UserDetails{
		FullName:        "Mr. Test",
		Email:           "mrtest@example.com",
		ListViewEnabled: 0,
		DarkModeEnabled: 1,
	}
	mockUsecase := new(mocks.UserUsecase)
	mockUsecase.On("GetUserDetails", mock.Anything, mock.AnythingOfType("int32")).Return(&mockUser, nil)

	handler := userHttp.UserHandler{
		UUseCase: mockUsecase,
	}

	t.Run("success", func(t *testing.T) {
		ctx, res := buildEchoAuthorizedRequest(t, echo.GET, endPoint, getToken(1), nil)
		handle := attachJWTMiddleware(handler.Me)

		assert.NoError(t, handle(ctx))
		assert.Equal(t, http.StatusOK, res.Code)

		var r response.Response
		assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &r))
		assert.True(t, true, r.Success)

		resultsMap := r.Results.(map[string]interface{})
		assert.Equal(t, mockUser.Email, resultsMap["email"])

		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		ctx, _ := buildEchoAuthorizedRequest(t, echo.GET, endPoint, "invalid_token", nil)
		handle := attachJWTMiddleware(handler.Me)

		assert.Error(t, handle(ctx))
		mockUsecase.AssertExpectations(t)
	})
}

// nolint:funlen
func TestUpdateSettings(t *testing.T) {
	endPoint := BaseURLV1 + "/me"

	hash, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	requestBody := updateSettings{
		OldPassword:     "",
		NewPassword:     "",
		ListViewEnabled: 0,
		DarkModeEnabled: 1,
	}
	mockUser := model.User{
		ID:              1,
		FullName:        "Mr. Test",
		Email:           "mrtest@example.com",
		Hash:            string(hash),
		IsActive:        1,
		IsTrashed:       0,
		ListViewEnabled: 0,
		DarkModeEnabled: 1,
	}

	mockUsecase := new(mocks.UserUsecase)
	mockUsecase.On("GetUser", mock.Anything, mock.AnythingOfType("int32")).Return(&mockUser, nil)
	mockUsecase.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	handler := userHttp.UserHandler{
		UUseCase: mockUsecase,
	}

	t.Run("update-settings", func(t *testing.T) {
		tempReq := requestBody

		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		ctx, res := buildEchoAuthorizedRequest(t, echo.POST, endPoint, getToken(1), strings.NewReader(string(j)))
		handle := attachJWTMiddleware(handler.UpdateSettings)

		assert.NoError(t, handle(ctx))
		assert.Equal(t, http.StatusOK, res.Code)

		var r response.Response
		assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &r))
		assert.True(t, true, r.Success)

		assert.Equal(t, "updated successfully", r.Message)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("update-password", func(t *testing.T) {
		tempReq := requestBody
		tempReq.OldPassword = "12345678"
		tempReq.NewPassword = "foobar-test"

		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		ctx, res := buildEchoAuthorizedRequest(t, echo.POST, endPoint, getToken(1), strings.NewReader(string(j)))
		handle := attachJWTMiddleware(handler.UpdateSettings)

		assert.NoError(t, handle(ctx))
		assert.Equal(t, http.StatusOK, res.Code)

		var r response.Response
		assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &r))
		assert.True(t, true, r.Success)

		assert.Equal(t, "updated successfully", r.Message)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("field-required", func(t *testing.T) {
		tempReq := requestBody
		tempReq.OldPassword = "12345678"

		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		ctx, res := buildEchoAuthorizedRequest(t, echo.POST, endPoint, getToken(1), strings.NewReader(string(j)))
		handle := attachJWTMiddleware(handler.UpdateSettings)

		assert.NoError(t, handle(ctx))
		assert.Equal(t, http.StatusBadRequest, res.Code)

		mockUsecase.AssertExpectations(t)
	})

	t.Run("field-required-2", func(t *testing.T) {
		tempReq := requestBody
		tempReq.NewPassword = "12345678"

		j, err := json.Marshal(tempReq)
		assert.NoError(t, err)
		ctx, res := buildEchoAuthorizedRequest(t, echo.POST, endPoint, getToken(1), strings.NewReader(string(j)))
		handle := attachJWTMiddleware(handler.UpdateSettings)

		assert.NoError(t, handle(ctx))
		assert.Equal(t, http.StatusBadRequest, res.Code)

		mockUsecase.AssertExpectations(t)
	})
}

func TestDeleteMe(t *testing.T) {
	endPoint := BaseURLV1 + "/me"

	hash, err := bcrypt.GenerateFromPassword([]byte("super_password"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	mockUser := model.User{
		ID:              1,
		FullName:        "Mr. Test",
		Email:           "mrtest@example.com",
		Hash:            string(hash),
		IsActive:        1,
		IsTrashed:       0,
		ListViewEnabled: 0,
		DarkModeEnabled: 1,
	}

	mockUsecase := new(mocks.UserUsecase)
	mockUsecase.On("GetUser", mock.Anything, mock.AnythingOfType("int32")).Return(&mockUser, nil)
	mockUsecase.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	handler := userHttp.UserHandler{
		UUseCase: mockUsecase,
	}

	t.Run("delete-account", func(t *testing.T) {
		ctx, res := buildEchoAuthorizedRequest(t, echo.DELETE, endPoint, getToken(1), nil)
		handle := attachJWTMiddleware(handler.DeleteMe)

		assert.NoError(t, handle(ctx))
		assert.Equal(t, http.StatusNoContent, res.Code)

		mockUsecase.AssertExpectations(t)
	})
}
