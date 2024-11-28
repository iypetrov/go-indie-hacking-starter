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
		Name     string
		Username string
		Password string
		Host     string
		Port     string
		SSL      string
	}
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg := &Config{}
	cfg.App.Addr = os.Getenv("APP_ADDR")
	cfg.App.Port = os.Getenv("APP_PORT")
	cfg.Database.Name = os.Getenv("APP_DB_NAME")
	cfg.Database.Username = os.Getenv("APP_DB_USERNAME")
	cfg.Database.Password = os.Getenv("APP_DB_PASSWORD")
	cfg.Database.Host = os.Getenv("APP_DB_HOST")
	cfg.Database.Port = os.Getenv("APP_DB_PORT")
	cfg.Database.SSL = os.Getenv("APP_DB_SSL")
	return cfg
}
