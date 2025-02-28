package processor

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection
var Ch *amqp.Channel

func ConnectToRabbitMQ(url string) error {
	var err error

	Conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	Ch, err = Conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func CloseConnection() {
	Ch.Close()
	Conn.Close()
}

func ConsumeMessage(body []byte) (ConsumedMessage, error) {
	var msg ConsumedMessage
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return ConsumedMessage{}, err
	}
	return msg, nil
}
