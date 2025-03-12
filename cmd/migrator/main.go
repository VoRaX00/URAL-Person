package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"os"
	"persons/internal/config"
)

type Config struct {
	DB             DBConfig `yaml:"db"`
	MigrationsPath string   `yaml:"migrations_path"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Password string `yaml:"-"`
	Username string `yaml:"username"`
	SSLMode  string `yaml:"ssl_mode"`
	IsDrop   bool   `yaml:"is_drop"`
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	path := config.FetchConfigPath()
	if path == "" {
		panic("path is empty")
	}
	cfg := config.MustConfig[Config](path)
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	db, err := sqlx.Open(cfg.DB.Driver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if cfg.DB.IsDrop {
		if err = goose.DownTo(db.DB, cfg.MigrationsPath, 0); err != nil {
			panic(err)
		}
	}

	if err = goose.Up(db.DB, cfg.MigrationsPath); err != nil {
		panic(err)
	}
}
