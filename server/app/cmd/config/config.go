package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BrianLusina/skillq/server/app/pkg/configs"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		configs.App   `yaml:"app"`
		configs.HTTP  `yaml:"http"`
		configs.Log   `yaml:"logger"`
		MongoDB       `yaml:"mongodb"`
		RabbitMQ      `yaml:"rabbitmq"`
		MinioConfig   `yaml:"minio"`
		ProductClient `yaml:"product_client"`
	}

	MongoDB struct {
		Host           string `env-required:"true" yaml:"host" env:"MONGODB_HOST"`
		Port           string `env-required:"true" yaml:"port" env:"MONGODB_PORT"`
		User           string `env-required:"true" yaml:"user" env:"MONGODB_USER"`
		Password       string `env-required:"true" yaml:"password" env:"MONGODB_PASSWORD"`
		RetryWrites    bool   `env-required:"false" yaml:"retryWrites" env:"MONGODB_RETRY_WRITES"`
		DatabaseName   string `env-required:"true" yaml:"database" env:"MONGODB_USER_DATABASE"`
		CollectionName string `env-required:"true" yaml:"collection" env:"MONGODB_USER_COLLECTION"`
	}

	RabbitMQ struct {
		Username string `env-required:"true" yaml:"username" env:"RABBITMQ_USERNAME"`
		Password string `env-required:"true" yaml:"password" env:"RABBITMQ_PASSWORD"`
		Host     string `env-required:"true" yaml:"host" env:"RABBITMQ_HOST"`
		Port     string `env-required:"true" yaml:"port" env:"RABBITMQ_PORT"`
		URL      string `env-required:"true" yaml:"url" env:"RABBITMQ_URL"`
	}

	MinioConfig struct {
		Endpoint        string `env-required:"true" yaml:"endpoint" env:"MINIO_ENDPOINT"`
		AccessKeyID     string `env-required:"true" yaml:"accessKeyId" env:"MINIO_ACCESS_KEY_ID"`
		SecretAccessKey string `env-required:"true" yaml:"secretAccessKey" env:"MINIO_SECRET_ACCESS_KEY"`
		UseSSL          bool   `env-required:"true" yaml:"useSSL" env:"MINIO_USE_SSL"`
		Token           string `env-required:"true" yaml:"token" env:"MINIO_TOKEN"`
	}

	ProductClient struct {
		URL string `env-required:"true" yaml:"url" env:"PRODUCT_CLIENT_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println("config path: " + dir)

	err = cleanenv.ReadConfig(dir+"/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
