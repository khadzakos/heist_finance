package processor

import (
	"log"
	"preprocessor/internal/config"
	"preprocessor/internal/storage"
	"strconv"
	"time"
)

type ConsumedMessage struct {
	Symbol   string `json:"s"`
	Price    string `json:"p"`
	Quantity string `json:"q"`
	Time     int64  `json:"T"`
}

func ProcessMessages(cfg config.PreprocessorConfig, db *storage.Storage) {
	msgs, err := Ch.Consume(
		cfg.Queue, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Ошибка подписки на очередь:", err)
	}

	for msg := range msgs {
		go func() {
			consumedMessage, err := ConsumeMessage(msg.Body)
			if err != nil {
				log.Printf("Ошибка обработки сообщения: %s\n", err)
				return
			}

			processedData := ProcessByExchange(cfg.Exchange, consumedMessage)
			db.SaveTradeData(processedData)
			log.Printf("Обработано (%s): %+v\n", cfg.Exchange, processedData)
		}()
	}
}

func ProcessByExchange(exchange string, msg ConsumedMessage) storage.MarketData {
	price, _ := strconv.ParseFloat(msg.Price, 64)
	return storage.MarketData{
		Exchange:  exchange,
		Timestamp: time.Now().Unix(),
		Ticker:    msg.Symbol,
		Price:     int64(price * 100000000), // 100000000 - 8 знаков после запятой в BIGINT
	}
}
