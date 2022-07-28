package mysql

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
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
	schemaPath = "file://../../infrastructure/db/migrations/mysql"
)

type MysqlRepositoryTestSuite struct {
	db *sql.DB
	suite.Suite
}

func (s *MysqlRepositoryTestSuite) SetupSuite() {
	if err := config.Load("./config.yml"); err != nil {
		logrus.WithError(err).Fatal("Failed to load config")
	}
	cfg := config.Get()
	connStr = fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db.Connect()
	s.db = db.GetClient()
}

func TestMysqlRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &MysqlRepositoryTestSuite{})
}

func (s *MysqlRepositoryTestSuite) SetupTest() {
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

func (s *MysqlRepositoryTestSuite) TearDownTest() {
	m, err := migrate.New(schemaPath, connStr)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), m.Down())
}
