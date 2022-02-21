package http

import (
	"librenote/app/response"
	"librenote/app/system/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// SystemHandler  represent the httphandler for system
type SystemHandler struct {
	Usecase usecase.SystemUsecase
}

// NewSystemHandler will initialize the system related endpoints
func NewSystemHandler(e *echo.Echo, us usecase.SystemUsecase) {
	handler := &SystemHandler{
		Usecase: us,
	}

	e.GET("/", handler.Root)
	e.GET("/h34l7h", handler.Health)
	e.GET("/api/v1/server-time", handler.ServerTime)
}

// Root will let you know, whoami
func (sh *SystemHandler) Root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "LibreNote API!"})
}

// Health will let you know the heart beats
func (sh *SystemHandler) Health(c echo.Context) error {
	resp, err := sh.Usecase.GetHealth()
	if err != nil {
		return c.JSON(response.RespondError(err))
	}
	return c.JSON(http.StatusOK, resp)
}

// ServerTime will let you know the server time
func (sh *SystemHandler) ServerTime(c echo.Context) error {
	resp := sh.Usecase.GetTime()
	return c.JSON(http.StatusOK, resp)
}
