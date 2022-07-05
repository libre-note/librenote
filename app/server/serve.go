package server

import (
	"context"
	"fmt"
	"librenote/app/model"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"librenote/app"
	systemDelivery "librenote/app/system/delivery/http"
	systemRepo "librenote/app/system/repository"
	systemUseCase "librenote/app/system/usecase"
	userDelivery "librenote/app/user/delivery/http"
	userMysqlRepo "librenote/app/user/repository/mysql"
	userPgsqlRepo "librenote/app/user/repository/pgsql"
	userSqliteRepo "librenote/app/user/repository/sqlite"
	userUseCase "librenote/app/user/usecase"
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
	contextTimeout := cfg.ContextTimeout * time.Second

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
	sysRepo := systemRepo.NewSystemRepository(dbClient)
	var uRepo model.UserRepository

	switch dbType {
	case "postgres":
		uRepo = userPgsqlRepo.NewPgsqlUserRepository(dbClient)
	case "mysql":
		uRepo = userMysqlRepo.NewMysqlUserRepository(dbClient)
	default:
		uRepo = userSqliteRepo.NewSqliteUserRepository(dbClient)
	}

	// use cases
	sysUseCase := systemUseCase.NewSystemUsecase(sysRepo)
	uUseCase := userUseCase.NewUserUsecase(uRepo, contextTimeout)

	// delivery
	systemDelivery.NewSystemHandler(e, sysUseCase)
	userDelivery.NewUserHandler(e, uUseCase)

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
