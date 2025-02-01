package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Connector struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

type Config struct {
	Connectors []Connector `yaml:"connectors"`
}

func LoadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Ошибка чтения config.yaml:", err)
	}
	var cfg Config
	yaml.Unmarshal(data, &cfg)
	return cfg
}
