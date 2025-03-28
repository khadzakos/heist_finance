package connectors

import (
	"context"

	"connector/internal/producer"
)

type ExchangeConnector interface {
	Connect(ctx context.Context) error
	SubscribeToMarketData(ctx context.Context, pub producer.MessageProducer) error
	// FetchHistoricalData(ctx context.Context, symbol string, period string, limit int) ([]HistoricalData, error)
}
