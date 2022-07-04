package db

import (
	"database/sql"
	"librenote/infrastructure/config"
	"os"

	"github.com/sirupsen/logrus"
)

type dbClient struct {
	DB *sql.DB
}

var dbc = &dbClient{}

func Connect() {
	dbType := config.Get().Database.Type
	switch dbType {
	case "postgres":
		if err := connectPG(); err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
	case "mysql":
		if err := connectMysql(); err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
	default:
		if err := connectSqlite(); err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
	}
}

func GetClient() *sql.DB {
	return dbc.DB
}

func Close() {
	_ = dbc.DB.Close()
}
