package processor

import (
	"preprocessor/internal/storage"
	"strconv"
	"time"
)

// TODO: реализовать обработку сообщений для всех сценариев

// ProcessByExchange - обрабатывает сообщение и возвращает структуру MarketData
func (w *Worker) ProcessByExchange(exchange string, msg ConsumedMessage) storage.MarketData {
	price, _ := strconv.ParseFloat(msg.Price, 64)
	price = price * 100000000
	return storage.MarketData{
		Exchange:  exchange,
		Symbol:    msg.Symbol,
		Price:     price,
		Bid:       price,
		Ask:       price,
		Timestamp: time.Now().Unix(),
	}
}
