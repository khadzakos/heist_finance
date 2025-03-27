package config

import (
	"os"
)

type Config struct {
	Exchange    string
	Queue       string
	RabbitMQURL string
}

func LoadConfig() Config {
	exchange := os.Getenv("EXCHANGE")
	queue := os.Getenv("QUEUE")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	return Config{
		Exchange:    exchange,
		Queue:       queue,
		RabbitMQURL: rabbitMQURL,
	}
}
