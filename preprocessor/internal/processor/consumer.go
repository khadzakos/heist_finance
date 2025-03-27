package processor

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

func (p *Processor) ConnectToRabbitMQ() error {
	conn, err := amqp.Dial(p.Cfg.RabbitMQ.URL)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	p.Conn = conn
	p.Ch = ch

	return nil
}

func (p *Processor) ConsumeMessage(body []byte) (ConsumedMessage, error) {
	var msg ConsumedMessage
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return ConsumedMessage{}, err
	}
	return msg, nil
}

func (p *Processor) CloseConnection() {
	p.Ch.Close()
	p.Conn.Close()
}
