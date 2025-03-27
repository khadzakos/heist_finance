package connectors

import (
	"context"

	"connector/internal/producer"
)

type ExchangeConnector interface {
	Connect(ctx context.Context) error
	ListenAndPublish(ctx context.Context, pub producer.MessageProducer) error
}
