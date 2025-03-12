package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"persons/internal/app"
	"persons/internal/app/server"
	"persons/internal/config"
	"persons/internal/handler"
	"persons/internal/service/person"
	"persons/internal/storage/postgres"
	personrepo "persons/internal/storage/postgres/person"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

const (
	configServer   = "./config/server.yml"
	configLogger   = "./config/logger.yml"
	configPostgres = "./config/postgres.yml"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	log := setupLogger(configLogger)

	pg := setupPostgres(configPostgres)
	personRepository := personrepo.NewRepository(pg)
	personService := person.NewService(log, personRepository, personRepository)
	h := handler.NewHandler(log, personService)

	srv := setupServer(configServer, h)
	application := app.NewApp(srv)

	log.Info("starting application")
	go application.MustStart()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutdown signal received")
	application.MustStop(context.Background())
	log.Info("shutdown complete")
}

func setupServer(configPath string, handler http.Handler) server.IServer {
	cfg := config.MustConfig[server.Config](configPath)
	srv := server.New(handler, cfg)
	return srv
}

func setupLogger(configPath string) *slog.Logger {
	cfg := config.MustConfig[config.Logger](configPath)

	var logger *slog.Logger
	switch cfg.Env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))
	}
	return logger
}

func setupPostgres(configPath string) *sqlx.DB {
	cfg := config.MustConfig[postgres.Config](configPath)
	cfg.Password = os.Getenv("DB_PASSWORD")

	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode),
	)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}
