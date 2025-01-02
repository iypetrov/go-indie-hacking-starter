package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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
	} else {
		cfg.AWS.Region = getSecretFromDockerSwarm("aws_region")
		cfg.AWS.AccessKeyID = getSecretFromDockerSwarm("aws_access_key_id")
		cfg.AWS.SecretAcessKey = getSecretFromDockerSwarm("aws_secret_access_key")
	}

	return cfg
}

func (cfg *Config) Load(ctx context.Context, awsCfg aws.Config) {
	if Profile == "local" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}

		cfg.App.Addr = os.Getenv("APP_ADDR")
		cfg.App.Port = "8080"
		cfg.Database.File = os.Getenv("APP_DB_FILE")
	} else {
		svc := secretsmanager.NewFromConfig(awsCfg)

		cfg.App.Addr = getSecretFromAWSSecretManager(ctx, svc, "go_indie_hacking_starter_addr")
		cfg.App.Port = "8080"
		cfg.Database.File = getSecretFromAWSSecretManager(ctx, svc, "go_indie_hacking_starter_db_file")
	}
}

func getSecretFromAWSSecretManager(ctx context.Context, svc *secretsmanager.Client, secretName string) string {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		panic(fmt.Errorf("can't find secret \"%s\" in aws secret manager", secretName))
	}

	return *result.SecretString
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
