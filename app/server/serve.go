package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"librenote/app"
	systemDelivery "librenote/app/system/delivery/http"
	systemRepo "librenote/app/system/repository"
	systemRepoPgsql "librenote/app/system/repository/pgsql"
	systemRepoSqlite "librenote/app/system/repository/sqlite"
	systemUseCase "librenote/app/system/usecase"
	"librenote/infrastructure/config"
	"librenote/infrastructure/db"
	"librenote/infrastructure/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Serve() {

	e := setupApiServer()

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// goroutine to handle shutdown
	go func() {
		// capture sigterm and other system call here
		<-sigCh
		logrus.Info("terminating...")

		// close db connection
		db.Close()

		// shutdown http server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = e.Shutdown(ctx)

	}()

	// start http server
	printBanner()
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.Get().App.Host, config.Get().App.Port)))
}

func setupApiServer() *echo.Echo {
	cfg := config.Get().App

	e := echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = cfg.ReadTimeout * time.Second
	e.Server.WriteTimeout = cfg.WriteTimeout * time.Second
	e.Server.IdleTimeout = cfg.IdleTimeout * time.Second

	// e.Validator = &validator.GenericValidator{}
	// fetch infra and routes
	if err := middlewares.Attach(e); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	db.Connect()
	dbClient := db.GetClient()
	dbType := config.Get().Database.Type

	// repository
	var sysRepo systemRepo.SystemRepository

	switch dbType {
	case "postgres":
		sysRepo = systemRepoPgsql.NewPgsqlSystemRepository(dbClient)
	default:
		sysRepo = systemRepoSqlite.NewSqliteSystemRepository(dbClient)
	}

	// use cases
	sysUseCase := systemUseCase.NewSystemUsecase(sysRepo)

	// delivery
	systemDelivery.NewSystemHandler(e, sysUseCase)

	return e
}

func printBanner() {
	log.SetFlags(0)
	log.Println("_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/")
	log.Println("_/                                                                                          _/")
	log.Println("_/                                                                                          _/")
	log.Println("_/    _/        _/  _/                            _/      _/              _/                _/")
	log.Println("_/     _/            _/_/_/    _/  _/_/    _/_/    _/_/    _/    _/_/    _/_/_/_/    _/_/   _/")
	log.Println("_/    _/        _/  _/    _/  _/_/      _/_/_/_/  _/  _/  _/  _/    _/    _/      _/_/_/_/  _/")
	log.Println("_/   _/        _/  _/    _/  _/        _/        _/    _/_/  _/    _/    _/      _/         _/")
	log.Println("_/  _/_/_/_/  _/  _/_/_/    _/          _/_/_/  _/      _/    _/_/        _/_/    _/_/_/    _/")
	log.Println("_/                                                                                          _/")
	log.Println("_/                                                                                          _/")
	log.Printf("_/                   Version: %-18s Build time: %-18s             _/", app.Version, app.BuildTime)
	log.Println("_/                                                                                          _/")
	log.Println("_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/")

}
