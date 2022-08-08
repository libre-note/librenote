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

type Server struct {
	ServerReady chan bool
}

func (s *Server) Serve() {
	cfg := config.Get().App

	db.Connect()

	if cfg.Env != "test" {
		defer db.Close()
	}

	e := setupAPIServer(cfg)

	go func() {
		printBanner()

		if err := e.Start(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
			logrus.Errorf(err.Error())
		}
	}()

	// server is ready now, so close the channel
	s.ServerReady <- true

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-sigCh
	logrus.Info("shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logrus.Fatalf("failed to gracefully shutdown the server: %s", err)
	}
}

func setupAPIServer(cfg config.AppConfig) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout
	contextTimeout := cfg.ContextTimeout

	if err := middlewares.Attach(e); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

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
