package main

import (
	"binance-connector/broker"
	"binance-connector/config"
	"binance-connector/wsclient"
)

func main() {
	cfg := config.LoadConfig()
	tradesChan := make(chan wsclient.TransactionData)
	client := wsclient.NewWebSocketClient(cfg.WsURL, cfg.Tickers, tradesChan)

	go client.ConnectWS()
	go broker.Produce(cfg, tradesChan)

	select {}
}
