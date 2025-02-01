package broker

import (
	"binance-connector/config"
	"binance-connector/wsclient"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func Produce(cfg config.Config, tradesChan chan wsclient.Transaction) {
	conn, err := amqp.Dial(config.GetRabbitMQConfig())
	if err != nil {
		log.Fatal("Ошибка подключения к RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Ошибка открытия канала:", err)
	}
	defer ch.Close()

	for {
		trade := <-tradesChan
		jsonTrade, err := json.Marshal(trade)
		if err != nil {
			log.Printf("Failed to marshal JSON: %v", err)
			continue
		}

		err = ch.Publish( // (autodeletion)
			"",
			cfg.QueueTopic,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonTrade,
			},
		)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		}
	}
}
