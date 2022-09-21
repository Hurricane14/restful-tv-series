package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"series/adapter/repository"
	"series/framework/database/postgres"
	"series/framework/handler/gorilla"
	"series/framework/validation/goplayground"
	"syscall"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Port int `env:"PORT"`
}

type DBConfig struct {
	Type       string `env:"TYPE"`
	Host       string `env:"HOST"`
	Port       int    `env:"PORT"`
	Database   string `env:"NAME"`
	User       string `env:"USER"`
	Password   string `env:"PASSWORD"`
	SSLEnabled string `env:"SSLMODE"`
}

type Config struct {
	Server ServerConfig
	DB     DBConfig `env-prefix:"DB_"`
}

func main() {
	var (
		config    Config
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

	handler = gorilla.NewHandler(
		repo,
		validator,
		10*time.Second,
	)
	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: handler,
	}

	log.Printf("Starting server at %s\n", server.Addr)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit

	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
