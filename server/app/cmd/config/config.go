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
		configs.App  `yaml:"app"`
		configs.HTTP `yaml:"http"`
		configs.Log  `yaml:"logger"`
		MongoDB      `yaml:"mongodb"`
		RabbitMQ     `yaml:"rabbitmq"`
		MinioConfig  `yaml:"minio"`
		EmailConfig  `yaml:"email"`
	}

	MongoDB struct {
		Host        string                       `env-description:"Mongo Database Host" yaml:"host" env:"MONGODB_HOST"`
		Port        string                       `env-description:"Mongo Database Port" yaml:"port" env:"MONGODB_PORT"`
		User        string                       `env-description:"Mongo User" yaml:"user" env:"MONGODB_USER"`
		Password    string                       `env-description:"Mongo Password" yaml:"password" env:"MONGODB_PASSWORD"`
		RetryWrites bool                         `env-description:"Mongo DB Retry writes" yaml:"retryWrites" env:"MONGODB_RETRY_WRITES"`
		Database    string                       `env-description:"Mongo DB database" yaml:"database" env:"MONGODB_USER_DATABASE"`
		Collections map[string]MongoDBCollection `env-description:"Mongo DB collection" yaml:"collections" env:"MONGODB_COLLECTIONS"`
	}

	MongoDBCollection struct {
		Name string `env-description:"Mongo DB USer collection" yaml:"user.name" env:"MONGODB_USER_COLLECTION"`
	}

	RabbitMQ struct {
		Username string `yaml:"username" env:"RABBITMQ_USERNAME"`
		Password string `yaml:"password" env:"RABBITMQ_PASSWORD"`
		Host     string `yaml:"host" env:"RABBITMQ_HOST"`
		Port     string `yaml:"port" env:"RABBITMQ_PORT"`
		URL      string `yaml:"url" env:"RABBITMQ_URL"`
	}

	MinioConfig struct {
		PublicUrl       string `yaml:"publicUrl" env:"MINIO_PUBLIC_URL"`
		Endpoint        string `yaml:"endpoint" env:"MINIO_ENDPOINT"`
		Host            string `yaml:"host" env:"MINIO_HOST"`
		Port            int    `yaml:"port" env:"MINIO_PORT"`
		AccessKeyID     string `yaml:"accessKeyId" env:"MINIO_ACCESS_KEY_ID"`
		SecretAccessKey string `yaml:"secretAccessKey" env:"MINIO_SECRET_ACCESS_KEY"`
		UseSSL          bool   `yaml:"useSSL" env:"MINIO_USE_SSL"`
		Token           string `yaml:"token" env:"MINIO_TOKEN"`
	}

	EmailConfig struct {
		Host     string `yaml:"host" env:"EMAIL_CLIENT_HOST"`
		Port     string `yaml:"port" env:"EMAIL_CLIENT_PORT"`
		Password string `yaml:"password" env:"EMAIL_CLIENT_PASSWORD"`
		From     string `yaml:"from" env:"EMAIL_CLIENT_FROM"`
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
