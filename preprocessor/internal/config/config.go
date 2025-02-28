package config

import (
	"os"
)

type RabbitMQConfig struct {
	URL string
}

type DBConfig struct {
	URL string
}

type PreprocessorConfig struct {
	Exchange string
	Queue    string
}

type Config struct {
	RabbitMQ      RabbitMQConfig
	Database      DBConfig
	Preprocessors PreprocessorConfig
}

func LoadConfig() Config {
	cfg := Config{
		RabbitMQ: RabbitMQConfig{
			URL: os.Getenv("RABBITMQ_URL"),
		},
		Database: DBConfig{
			URL: os.Getenv("DATABASE_URL"),
		},
		Preprocessors: PreprocessorConfig{
			Exchange: os.Getenv("EXCHANGE"),
			Queue:    os.Getenv("QUEUE"),
		},
	}

	return cfg
}
