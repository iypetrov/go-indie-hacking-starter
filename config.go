package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AWS struct {
		Region         string
		AccessKeyID    string
		SecretAcessKey string
	}

	App struct {
		Domain string
		Port   string
	}

	Database struct {
		File string
	}
}

func NewConfig() *Config {
	cfg := &Config{}

	if Profile == "local" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}

		cfg.App.Domain = os.Getenv("APP_DOMAIN")
		cfg.App.Port = os.Getenv("APP_PORT")
		cfg.Database.File = os.Getenv("DB_FILE")
		cfg.AWS.Region = os.Getenv("AWS_REGION")
		cfg.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
		cfg.AWS.SecretAcessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	} else {
		// TODO: Rename them later to go_indie_hacking_starter_*
		cfg.AWS.Region = getSecretFromDockerSwarm("blog_aws_region")
		cfg.AWS.AccessKeyID = getSecretFromDockerSwarm("blog_aws_access_key_id")
		cfg.AWS.SecretAcessKey = getSecretFromDockerSwarm("blog_aws_secret_access_key")
		cfg.App.Domain = getSecretFromDockerSwarm("blog_domain")
		cfg.App.Port = getSecretFromDockerSwarm("blog_port")
		cfg.Database.File = getSecretFromDockerSwarm("blog_db_file")
	}

	return cfg
}

func getSecretFromDockerSwarm(secretName string) string {
	secretFile, err := os.Open("/run/secrets/" + secretName)
	if err != nil {
		panic(fmt.Errorf("can't open secret \"%s\" in docker swarm", secretName))
	}
	defer secretFile.Close()

	secretContent, err := os.ReadFile("/run/secrets/" + secretName)
	if err != nil {
		panic(fmt.Errorf("can't find secret \"%s\" in docker swarm", secretName))
	}

	return string(secretContent)
}
