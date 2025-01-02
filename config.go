package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/joho/godotenv"
)

type Config struct {
	AWS struct {
		Region         string
		AccessKeyID    string
		SecretAcessKey string
	}

	App struct {
		Addr string
		Port string
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

		cfg.AWS.Region = os.Getenv("AWS_REGION")
		cfg.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
		cfg.AWS.SecretAcessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		cfg.App.Addr = os.Getenv("APP_ADDR")
		cfg.App.Port = "8080"
		cfg.Database.File = os.Getenv("APP_DB_FILE")
	} else {
		cfg.AWS.Region = os.Getenv("AWS_REGION")
		cfg.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
		cfg.AWS.SecretAcessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(cfg.AWS.Region),
		})
		if err != nil {
			panic(err)
		}
		sm := secretsmanager.New(sess)

		cfg.App.Addr = getSecretFromAWSSecretManager(sm, "go_indie_hacking_starter_addr")
		cfg.App.Port = "8080"
		cfg.Database.File = getSecretFromAWSSecretManager(sm, "go_indie_hacking_starter_db_file")
	}

	return cfg
}

func getSecretFromAWSSecretManager(sm *secretsmanager.SecretsManager, secretName string) string {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}
	result, err := sm.GetSecretValue(input)
	if err != nil {
		panic(fmt.Errorf("can't find secret \"%s\" in aws secret manager", secretName))
	}

	return *result.SecretString
}
