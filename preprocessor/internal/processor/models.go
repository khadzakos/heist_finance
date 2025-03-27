package processor

import (
	"preprocessor/internal/config"
	"preprocessor/internal/storage"

	"github.com/streadway/amqp"
)

type Processor struct {
	Cfg  *config.Config
	Db   *storage.Storage
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

type ConsumedMessage struct {
	Symbol   string `json:"s"`
	Price    string `json:"p"`
	Quantity string `json:"q"`
	Time     int64  `json:"T"`
}

type Worker struct {
	Id        int
	Jobs      <-chan amqp.Delivery
	Db        *storage.Storage
	Processor *Processor
}
