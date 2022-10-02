package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"series/adapter/logger"
	"series/adapter/repository"
	"series/framework/database/postgres"
	"series/framework/handler/gorilla"
	"series/framework/logging/logrus"
	"series/framework/validation/goplayground"
	"syscall"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port int `env:"PORT"`
	}
	DB struct {
		Type       string `env:"TYPE"`
		Host       string `env:"HOST"`
		Port       int    `env:"PORT"`
		Database   string `env:"NAME"`
		User       string `env:"USER"`
		Password   string `env:"PASSWORD"`
		SSLEnabled string `env:"SSLMODE"`
	} `env-prefix:"DB_"`
	Logger struct {
		Level string `env:"LVL"`
		File  string `env:"FILE" env-default:"stderr"`
	} `env-prefix:"LOG_"`
}

func main() {
	var (
		config    Config
		logger    logger.Logger
		repo      repository.Repository
		validator = goplayground.NewValidator()
		handler   http.Handler
		server    *http.Server
	)

	if err := cleanenv.ReadEnv(&config); err != nil {
		panic(err)
	}

	switch config.DB.Type {
	case "postgres":

		cfg := postgres.NewConfig().
			WithHost(config.DB.Host).
			WithPort(config.DB.Port).
			WithDatabase(config.DB.Database).
			WithUser(config.DB.User).
			WithPassword(config.DB.Password).
			WithSSLMode(config.DB.SSLEnabled)

		db, err := postgres.NewDB(cfg)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		repo = db

	default:
		panic("Unknown db type")
	}

	{
		var err error
		var logFile = os.Stderr
		if filename := config.Logger.File; filename != "stderr" {
			logFile, err = os.OpenFile(
				filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600,
			)
			if err != nil {
				panic(err)
			}
			defer logFile.Close()
		}

		logger, err = logrus.New(config.Logger.Level, logFile)
		if err != nil {
			panic(err)
		}
	}

	handler = gorilla.NewHandler(
		repo,
		logger,
		validator,
		10*time.Second,
	)
	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: handler,
	}

	logger.Infof("Starting server at %s", server.Addr)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit

	logger.Infof("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
