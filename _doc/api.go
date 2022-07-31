// @title LibreNote API
// @version 1.0
// @description This is a LibreNote API server
// @x-logo {"url": "https://avatars.githubusercontent.com/u/99524695?s=200&v=4", "backgroundColor": "#FFFFFF", "altText": "LibreNote Logo"}

// @contact.name API Support
// @contact.url https://github.com/libre-note/librenote
// @contact.email dev@hrshadhin.me

// @license.name MIT
// @license.url https://github.com/libre-note/librenote/blob/master/LICENSE.md

// @host localhost:8000
// @BasePath /
// @schemes https
// @accept json

// Package _doc provides basic API structs for the REST services
package _doc

// RootResponse is a data structure for the / endpoint
type RootResponse struct {
	// service details message
	Message string `json:"message"`
}

// HealthResponse is a data structure for the /h34l7h endpoint
type HealthResponse struct {
	// health details
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TimeResp is a data structure for the /time endpoint
type TimeResp struct {
	// current server time
	CurrentTimeUnix int64 `json:"current_time_unix"`
}

type successResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type failedResponse struct {
	Success bool   `json:"success" default:"false"`
	Message string `json:"message"`
}

type loginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type successResponseData struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Results interface{} `json:"result"`
}

type registrationReq struct {
	FullName string `json:"full_name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type loginReq struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// Root will let you know, whoami
// @Summary Service Details
// @Description will let you know, whoami
// @Tags system
// @Produce	json
// @Success	200	{object} RootResponse "return service details"
// @Failure	500	{object} failedResponse
// @Router / [get]
func Root() {}

// Health will let you know the heart beats
// @Summary Service Health
// @Description will let you know the heart beats
// @Tags system
// @Produce	json
// @Success	200	{object} HealthResponse "return service health"
// @Failure	500	{object} failedResponse
// @Router /h34l7h [get]
func Health() {}

// ServerTime will let you know the server time
// @Summary Server time
// @Description will let you know the server time
// @Tags system
// @Produce	json
// @Success	200	{object} TimeResp "return server current time"
// @Failure	500	{object} failedResponse
// @Router /time [get]
func ServerTime() {}

// Registration
// @Summary Registration
// @Description user registration endpoint
// @Tags user
// @Accept json
// @Param payload body registrationReq false "Registration Payload"
// @Produce	json
// @Success	200	{object} successResponse
// @Failure	400,422,500	{object} failedResponse
// @Router /api/v1/registration [post]
func Registration() {}

// Login
// @Summary Login
// @Description user login endpoint
// @Tags user
// @Accept json
// @Param payload body loginReq false "Login Payload"
// @Produce	json
// @Success	200	{object} loginResponse
// @Failure	400,401,422,500	{object} failedResponse
// @Router /api/v1/login [post]
func Login() {}

// Me
// @Summary Me
// @Description user details endpoint
// @Tags user
// @Param Authorization header string true "Bearer {Token}"
// @Produce	json
// @Success	200	{object} model.UserDetails
// @Failure	401,403,404,500	{object} failedResponse
// @Router /api/v1/me [get]
func Me() {}
