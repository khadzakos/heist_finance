package processor

import (
	"encoding/json"
	"fmt"

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

func (p *Processor) ConsumeMessage(body []byte) (GenericMessage, error) {
	var msg GenericMessage
	switch p.Cfg.Preprocessor.Exchange {
	case "binance":
		var msgBinance BinanceMarketData
		err := json.Unmarshal(body, &msgBinance)
		if err != nil {
			return nil, err
		}
		msg = msgBinance
	case "bybit":
		var msgBybit BybitMarketData
		err := json.Unmarshal(body, &msgBybit)
		if err != nil {
			return nil, err
		}
		msg = msgBybit
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", p.Cfg.Preprocessor.Exchange)
	}
	return msg, nil
}

func (p *Processor) CloseConnection() {
	p.Ch.Close()
	p.Conn.Close()
}
