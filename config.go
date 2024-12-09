package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Addr string
		Port string
	}

	Database struct {
		File string
	}
}

func NewConfig() *Config {
	if Profile == "local" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	cfg := &Config{}
	cfg.App.Addr = os.Getenv("APP_ADDR")
	cfg.App.Port = os.Getenv("APP_PORT")
	cfg.Database.File = os.Getenv("APP_DB_FILE")
	return cfg
}
