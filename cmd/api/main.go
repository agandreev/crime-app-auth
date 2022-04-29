package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/agandreev/crime-app-auth/internal/handlers"
	"github.com/agandreev/crime-app-auth/internal/repository"
	"github.com/agandreev/crime-app-auth/internal/service"
	"github.com/agandreev/crime-app-auth/pkg/config"
	"github.com/agandreev/crime-app-auth/pkg/database"
	"github.com/agandreev/crime-app-auth/pkg/server"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	configPath = "./config/config.yaml"
)

// @title Crime app auth
// @version 1.0
// @description Crime app auth provides authentication for crime-app microservices.

// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()
	log := logrus.New()

	cfg, err := config.Read(configPath)
	if err != nil {
		log.Fatal("can't load config, err: %w", err)
	}

	db, err := database.Connect(ctx, cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}

	// create layers
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, cfg.JWTConfig)

	e := echo.New()
	handlers.InitCommonRoutes(e)
	userHandler := handlers.NewUserHandler(userService, log)
	userHandler.InitRoutes(e, cfg.CtxTimeout)

	srv := server.NewServer(cfg.SrvConfig, e)

	// gracefully shutdown
	go func() {
		if err = srv.Run(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("running server is failed, err: %q", err)
			}
		}
	}()

	log.Info("server is running...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("server is shutting down...")

	if err = srv.Stop(ctx); err != nil {
		log.Fatalf("graceful shutdown is broken, err: %q", err)
	}
	log.Info("srv was gracefully shut down.")
}
