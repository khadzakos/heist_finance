package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Connector struct {
	Name        string `yaml:"name"`
	Image       string `yaml:"image"`
	Exchange    string `yaml:"exchange"`
	Queue       string `yaml:"queue"`
	RabbitMQURL string
}

type Preprocessor struct {
	Name        string `yaml:"name"`
	Exchange    string `yaml:"exchange"`
	Image       string `yaml:"image"`
	Queue       string `yaml:"queue"`
	RabbitMQURL string
	DatabaseURL string
}

type Environment struct {
	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string
	RABBITMQ_HOST     string
	DATABASE_HOST     string
	DATABASE_USER     string
	DATABASE_PASSWORD string
	DATABASE_NAME     string
	DATABASE_PORT     string
}

type Config struct {
	Network       string         `yaml:"network"`
	Connectors    []Connector    `yaml:"connectors"`
	Preprocessors []Preprocessor `yaml:"preprocessors"`
}

func LoadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Ошибка чтения config.yaml:", err)
	}
	var cfg Config
	yaml.Unmarshal(data, &cfg)

	env := Environment{
		RABBITMQ_USER:     os.Getenv("RABBITMQ_USER"),
		RABBITMQ_PASSWORD: os.Getenv("RABBITMQ_PASSWORD"),
		RABBITMQ_HOST:     os.Getenv("RABBITMQ_HOST"),
		DATABASE_HOST:     os.Getenv("DATABASE_HOST"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
		DATABASE_NAME:     os.Getenv("DATABASE_NAME"),
		DATABASE_PORT:     os.Getenv("DATABASE_PORT"),
	}

	rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:5672", env.RABBITMQ_USER, env.RABBITMQ_PASSWORD, env.RABBITMQ_HOST)
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", env.DATABASE_USER, env.DATABASE_PASSWORD, env.DATABASE_HOST, env.DATABASE_PORT, env.DATABASE_NAME)
	for i := range cfg.Connectors {
		cfg.Connectors[i].RabbitMQURL = rabbitMQURL
	}

	for i := range cfg.Preprocessors {
		cfg.Preprocessors[i].RabbitMQURL = rabbitMQURL
		cfg.Preprocessors[i].DatabaseURL = databaseURL
	}

	return cfg
}

func GetDatabaseURL() string {
	env := Environment{
		DATABASE_HOST:     os.Getenv("DATABASE_HOST"),
		DATABASE_USER:     os.Getenv("DATABASE_USER"),
		DATABASE_PASSWORD: os.Getenv("DATABASE_PASSWORD"),
		DATABASE_NAME:     os.Getenv("DATABASE_NAME"),
		DATABASE_PORT:     os.Getenv("DATABASE_PORT"),
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", env.DATABASE_USER, env.DATABASE_PASSWORD, env.DATABASE_HOST, env.DATABASE_PORT, env.DATABASE_NAME)

	return databaseURL
}
