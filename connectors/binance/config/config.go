package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	WsURL   string
	Tickers []string
	DBURL   string
}

func LoadConfig() Config {
	wsURL := os.Getenv("WS_URL")
	tickers := os.Getenv("TICKERS")

	if wsURL == "" || tickers == "" {
		log.Fatal("Binance: Отсутствуют ENV-переменные: WS_URL, TICKERS")
	}

	return Config{
		WsURL:   wsURL,
		Tickers: strings.Split(tickers, ","),
	}
}
