package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	WsURL      string   `yaml:"ws_url"`
	Tickers    []string `yaml:"tickers"`
	QueueTopic string   `yaml:"queue_topic"`
}

func LoadConfig() Config {
	data, err := os.ReadFile("") // TODO: fix (config.yaml)
	if err != nil {
		log.Fatal("Ошибка чтения конфигурации:", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("Ошибка парсинга конфигурации:", err)
	}
	return cfg
}

func GetRabbitMQConfig() string {
	err := godotenv.Load("") // TODO: fix (secret.env)
	if err != nil {
		log.Fatal("Ошибка загрузки переменных окружения:", err)
	}

	login := os.Getenv("RABBITMQ_LOGIN")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	connURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", login, password, host, port)
	return connURL
}
