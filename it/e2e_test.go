package it

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"librenote/app/model"
	"librenote/app/response"
	"librenote/app/server"
	repo "librenote/app/user/repository/sqlite"
	"librenote/infrastructure/config"
	"librenote/infrastructure/db"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type e2eTestSuite struct {
	suite.Suite
	db          *sql.DB
	dbMigration *migrate.Migrate
	apiBaseURL  string
}

const loginJSON = `{"email": "mrtest3@example.com", "password":"12345678"}`

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	s.Require().NoError(config.Load("./config.yml"))

	cfg := config.Get()
	connectionStr := fmt.Sprintf("sqlite3://%s/%s.db", cfg.App.DataPath, cfg.Database.Name)

	var err error

	s.dbMigration, err = migrate.New("file://../infrastructure/db/migrations/sqlite", connectionStr)
	s.Require().NoError(err)

	if err := s.dbMigration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err)
	}

	serverReady := make(chan bool)
	httpServer := server.Server{ServerReady: serverReady}

	go httpServer.Serve()

	// wait until api server is start
	<-serverReady

	s.db = db.GetClient()
	s.apiBaseURL = fmt.Sprintf("http://localhost:%d/api/v1", cfg.App.Port)
}

func (s *e2eTestSuite) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	_ = p.Signal(syscall.SIGINT)
}

func (s *e2eTestSuite) SetupTest() {
	if err := s.dbMigration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err)
	}
}

func (s *e2eTestSuite) TearDownTest() {
	s.NoError(s.dbMigration.Down())
}

func (s *e2eTestSuite) Test_EndToEnd_RegisterUser() {
	reqStr := `{"full_name":"Mr. Test", "email": "test01@example.com", "password":"12345678"}`
	req, err := http.NewRequest(echo.POST, s.apiBaseURL+"/registration", strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)

	byteBody, err := io.ReadAll(res.Body)
	s.NoError(err)

	s.Equal(`{"success":true,"message":"registration successful"}`, strings.Trim(string(byteBody), "\n"))

	_ = res.Body.Close()
}

func (s *e2eTestSuite) Test_EndToEnd_LoginUser() {
	s.createUser(3)

	reqStr := `{"email": "mrtest1@example.com", "password":"12345678"}`
	_ = s.doLogin(reqStr)
}

func (s *e2eTestSuite) doLogin(payload string) string {
	req, err := http.NewRequest(echo.POST, s.apiBaseURL+"/login", strings.NewReader(payload))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	byteBody, err := io.ReadAll(res.Body)
	s.NoError(err)

	_ = res.Body.Close()

	var r response.Response

	s.NoError(json.Unmarshal(byteBody, &r))

	s.True(true, r.Success)
	s.Equal("Login successful", r.Message)

	return r.Token
}

func (s *e2eTestSuite) createUser(howMany int) {
	for i := 1; i <= howMany; i++ {
		nowTime := time.Now().UTC().Format("2006-01-02 15:04:05")
		newUser := &model.User{
			FullName:  fmt.Sprintf("Mr. Test %d", i),
			Email:     fmt.Sprintf("mrtest%d@example.com", i),
			Hash:      "12345678",
			IsActive:  1,
			IsTrashed: 0,
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Hash), bcrypt.MinCost)

		s.Assert().NoError(err)

		newUser.Hash = string(hash)
		r := repo.NewSqliteUserRepository(s.db)

		s.Assert().NoError(r.CreateUser(context.Background(), newUser))
	}
}

func (s *e2eTestSuite) Test_EndToEnd_Me() {
	s.createUser(5)

	reqStr := `{"email": "mrtest3@example.com", "password":"12345678"}`
	token := s.doLogin(reqStr)

	req, err := http.NewRequest(echo.GET, s.apiBaseURL+"/me", nil)
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	byteBody, err := io.ReadAll(res.Body)

	s.NoError(err)

	_ = res.Body.Close()

	var r response.Response

	s.NoError(json.Unmarshal(byteBody, &r))

	s.True(true, r.Success)

	resultsMap := r.Results.(map[string]interface{})
	s.Equal("mrtest3@example.com", resultsMap["email"])
}

func (s *e2eTestSuite) Test_EndToEnd_UpdateSettings() {
	s.createUser(5)

	reqStr := `{"email": "mrtest3@example.com", "password":"12345678"}`
	token := s.doLogin(reqStr)

	reqStr = `{"list_view_enabled":1, "dark_mode_enabled":0}`
	req, err := http.NewRequest(echo.POST, s.apiBaseURL+"/me", strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	byteBody, err := io.ReadAll(res.Body)

	s.NoError(err)

	_ = res.Body.Close()

	var r response.Response

	s.NoError(json.Unmarshal(byteBody, &r))

	s.True(true, r.Success)
	s.Equal("updated successfully", r.Message)
}

func (s *e2eTestSuite) Test_EndToEnd_UpdatePassword() {
	s.createUser(5)

	token := s.doLogin(loginJSON)

	reqStr := `{"old_password":"12345678", "new_password":"super_passwrod_updated", "list_view_enabled":1,
"dark_mode_enabled":0}`
	req, err := http.NewRequest(echo.POST, s.apiBaseURL+"/me", strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	byteBody, err := io.ReadAll(res.Body)

	s.NoError(err)

	_ = res.Body.Close()

	var r response.Response

	s.NoError(json.Unmarshal(byteBody, &r))

	s.True(true, r.Success)
	s.Equal("updated successfully", r.Message)
}

func (s *e2eTestSuite) Test_EndToEnd_DeleteAccount() {
	s.createUser(3)

	token := s.doLogin(loginJSON)

	req, err := http.NewRequest(echo.DELETE, s.apiBaseURL+"/me", nil)
	s.NoError(err)

	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	_ = res.Body.Close()

	s.Equal(http.StatusNoContent, res.StatusCode)
}
