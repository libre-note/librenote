package cmd

import (
	"database/sql"
	"fmt"

	"librenote/infrastructure/config"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var schemaPath string

func migrateCommand() *cobra.Command {
	c := &cobra.Command{
		Use:              "migrate",
		Short:            "migrate database",
		Long:             `migrate database like(postgres, mysql, sqlite)`,
		TraverseChildren: true,
	}
	c.PersistentFlags().StringVarP(&schemaPath, "path", "p", "", "migration file path")
	_ = c.MarkPersistentFlagRequired("path")

	c.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "migration up",
		Long:  `migration up`,
		Run: func(cmd *cobra.Command, args []string) {
			step := 0
			if len(args) > 0 {
				if _, err := strconv.Atoi(args[0]); err != nil {
					logrus.Printf("%q migration step should be a number\n", args[0])
					os.Exit(1)
				}
			}
			err := migrateDatabase("up", step)
			if err != nil {
				logrus.Printf("%v", err)
				os.Exit(1)
			}
		},
	})

	c.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "migration down",
		Long:  `migration down`,
		Run: func(cmd *cobra.Command, args []string) {
			step := 0
			if len(args) > 0 {
				if _, err := strconv.Atoi(args[0]); err != nil {
					logrus.Printf("%q migration step should be a number\n", args[0])
					os.Exit(1)
				}
			}
			err := migrateDatabase("down", step)
			if err != nil {
				logrus.Printf("%v", err)
				os.Exit(1)
			}
		},
	})

	return c
}

func migrateDatabase(state string, step int) error {
	cfg := config.Get()
	dbType := cfg.Database.Type
	driverName := "sqlite3"
	dbURL := fmt.Sprintf("%s/%s.db", cfg.App.DataPath, cfg.Database.Name)

	switch dbType {
	case "postgres":
		driverName = "pgx"
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SslMode,
		)
	case "mysql":
		driverName = "mysql"
	}

	db, err := sql.Open(driverName, dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	var instance database.Driver
	if dbType == "postgres" {
		instance, err = pgx.WithInstance(db, &pgx.Config{})
		if err != nil {
			return err
		}

	} else {
		instance, err = sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			return err
		}
	}

	fSrc, err := (&file.File{}).Open(schemaPath)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("file", fSrc, driverName, instance)
	if err != nil {
		return err
	}

	if state == "up" {
		if step > 0 {
			if err := m.Steps(step); err != nil {
				return err
			}
		} else {
			if err := m.Up(); err != nil {
				return err
			}
		}
	}

	if state == "down" {
		if step > 0 {
			if err := m.Steps(-step); err != nil {
				return err
			}
		} else {
			if err := m.Down(); err != nil {
				return err
			}
		}
	}

	return nil
}
