package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() *Config {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
	)

	return &Config{
		DatabaseURL: dbURL,
	}
}
