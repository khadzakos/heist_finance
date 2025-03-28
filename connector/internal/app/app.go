package app

import (
	"context"
	"log"

	"connector/internal/config"
	"connector/internal/producer"

	"connector/internal/connectors"
	"connector/internal/connectors/binance"
	"connector/internal/connectors/bybit"
)

func Run() {
	cfg := config.LoadConfig()

	log.Println(cfg)

	var connector connectors.ExchangeConnector
	switch cfg.Exchange {
	case "binance":
		connector = binance.NewConnector()
	case "bybit":
		connector = bybit.NewConnector()
	// case "okx":
	// 	connector = okx.NewConnector()
	default:
		log.Fatalf("unsupported exchange: %s", cfg.Exchange)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pub, err := producer.NewRabbitProducer(cfg.RabbitMQURL, cfg.Queue)
	if err != nil {
		log.Fatalf("create producer: %v", err)
	}
	defer pub.Close()

	if err := connector.Connect(ctx); err != nil {
		log.Fatalf("connect: %v", err)
	}
	if err := connector.SubscribeToMarketData(ctx, pub); err != nil {
		log.Fatalf("listen & publish: %v", err)
	}
}
