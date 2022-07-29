package it

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"librenote/app/model"
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
)

type e2eTestSuite struct {
	suite.Suite
	db          *sql.DB
	dbMigration *migrate.Migrate
	apiBaseURL  string
}

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
	if err := s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
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
	p.Signal(syscall.SIGINT)
}

func (s *e2eTestSuite) SetupTest() {
	if err := s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
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
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	s.Equal(`{"success":true,"message":"registration successful"}`, strings.Trim(string(byteBody), "\n"))
	_ = response.Body.Close()
}

func (s *e2eTestSuite) Test_EndToEnd_LoginUser() {
	nowTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	newUser := &model.User{
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
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

	reqStr := `{"email": "mrtest@example.com", "password":"12345678"}`
	req, err := http.NewRequest(echo.POST, s.apiBaseURL+"/login", strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)
	s.True(strings.Contains(strings.Trim(string(byteBody), "\n"), "Login successful"))
	_ = response.Body.Close()
}
