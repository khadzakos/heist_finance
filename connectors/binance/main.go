package main

import (
	"log"

	"binance-connector/config"
	"binance-connector/wsclient"
)

func main() {
	cfg := config.LoadConfig()
	client := wsclient.NewWebSocketClient(cfg.WsURL, cfg.Tickers)

	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}

	client.Listen()
}
