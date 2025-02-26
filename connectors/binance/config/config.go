package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WsURL   string   `yaml:"ws_url"`
	Tickers []string `yaml:"tickers"`
	Queue   string
}

func LoadConfig() Config {
	queue := os.Getenv("QUEUE")

	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatal("Ошибка чтения конфигурации:", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("Ошибка парсинга конфигурации:", err)
	}
	cfg.Queue = queue
	return cfg
}

func GetRabbitMQConfig() string {
	login := os.Getenv("RABBITMQ_URL")
	return login
}
