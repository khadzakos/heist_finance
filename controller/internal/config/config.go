package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConnectorConfig struct {
	Name    string   `yaml:"name"`
	Image   string   `yaml:"image"`
	WsUrl   string   `yaml:"ws_url"`
	Tickers []string `yaml:"tickers"`
}

type Config struct {
	Connectors []ConnectorConfig `yaml:"connectors"`
}

func LoadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Ошибка чтения config.yaml: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Ошибка парсинга YAML: %v", err)
	}

	return cfg
}
