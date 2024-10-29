//go:build local
// +build local

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var cfg *Config

type Config struct {
	App struct {
		Profile string
		Version string
		Addr    string
		Port    string
	}
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg = &Config{}
	cfg.App.Profile = string(Local)
	cfg.App.Version = os.Getenv("APP_VERSION")
	cfg.App.Addr = os.Getenv("APP_ADDR")
	cfg.App.Port = os.Getenv("APP_PORT")
}

func Get() *Config {
	return cfg
}

func (c *Config) BaseWebUrl() string {
	protocol := "https://"
	basePath := c.App.Addr

	if c.App.Profile == string(Local) {
		protocol = "http://"
		basePath = fmt.Sprintf("%s:%s", c.App.Addr, c.App.Port)
	}

	return fmt.Sprintf("%s%s", protocol, basePath)
}

func (c *Config) ViewPrefix() string {
	return "/p"
}

func (c *Config) PublicViewPrefix() string {
	return "/public"
}

func (c *Config) ClientViewPrefix() string {
	return "/client"
}

func (c *Config) AdminViewPrefix() string {
	return "/admin"
}

func (c *Config) PublicApiPrefix() string {
	return fmt.Sprintf("/public/v%s", c.App.Version)
}

func (c *Config) ClientApiPrefix() string {
	return fmt.Sprintf("/api/v%s", c.App.Version)
}

func (c *Config) AdminApiPrefix() string {
	return fmt.Sprintf("/admin/v%s", c.App.Version)
}
