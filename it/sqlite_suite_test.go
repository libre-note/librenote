package it_test

import (
	"database/sql"
	"errors"
	"fmt"
	"librenote/infrastructure/config"
	"librenote/infrastructure/db"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	connStr    string
	schemaPath = "file://../infrastructure/db/migrations/sqlite"
)

type SqliteRepositoryTestSuite struct {
	db *sql.DB
	suite.Suite
}

func (s *SqliteRepositoryTestSuite) SetupSuite() {
	if err := config.Load("./config.yml"); err != nil {
		logrus.WithError(err).Fatal("Failed to load config")
	}

	cfg := config.Get()
	connStr = fmt.Sprintf("sqlite3://%s/%s.db", cfg.App.DataPath, cfg.Database.Name)

	db.Connect()
	s.db = db.GetClient()
}

func TestSqliteRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &SqliteRepositoryTestSuite{})
}

func (s *SqliteRepositoryTestSuite) SetupTest() {
	m, err := migrate.New(schemaPath, connStr)
	assert.NoError(s.T(), err)

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			// just ignore
			return
		}

		panic(err)
	}
}

func (s *SqliteRepositoryTestSuite) TearDownTest() {
	m, err := migrate.New(schemaPath, connStr)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), m.Down())
}
