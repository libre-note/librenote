package http_test

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"librenote/app/model/mocks"
	userHttp "librenote/app/user/delivery/http"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var BaseURLV1 = "/api/v1"

type registrationReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestRegistration(t *testing.T) {
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
		c, rec := buildEchoRequest(t, endPoint, strings.NewReader(string(j)))

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
		c, rec := buildEchoRequest(t, endPoint, strings.NewReader(string(j)))

		handler := userHttp.UserHandler{
			UUseCase: mockUsecase,
		}
		err = handler.Registration(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func buildEchoRequest(t *testing.T, path string, payload io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req, err := http.NewRequest(echo.POST, path, payload)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)

	return c, rec
}
