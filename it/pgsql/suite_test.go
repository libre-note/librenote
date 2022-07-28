package pgsql

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"librenote/infrastructure/config"
	"librenote/infrastructure/db"
	"testing"
)

var (
	connStr    string
	schemaPath = "file://../../infrastructure/db/migrations/pgsql"
)

type PgsqlRepositoryTestSuite struct {
	db *sql.DB
	suite.Suite
}

func (s *PgsqlRepositoryTestSuite) SetupSuite() {
	if err := config.Load("./config.yml"); err != nil {
		logrus.WithError(err).Fatal("Failed to load config")
	}
	cfg := config.Get()
	connStr = fmt.Sprintf("pgx://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db.Connect()
	s.db = db.GetClient()
}

func TestPgsqlRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &PgsqlRepositoryTestSuite{})
}

func (s *PgsqlRepositoryTestSuite) SetupTest() {
	m, err := migrate.New(schemaPath, connStr)
	assert.NoError(s.T(), err)

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			// just ignore
			return
		}

		panic(err)
	}
}

func (s *PgsqlRepositoryTestSuite) TearDownTest() {
	m, err := migrate.New(schemaPath, connStr)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), m.Down())
}
