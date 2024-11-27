package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Addr    string
		Port    string
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
	return cfg
}
