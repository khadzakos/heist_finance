package processor

import (
	"log"
	"preprocessor/internal/storage"
	"strconv"
)

// TODO: реализовать обработку сообщений для всех сценариев

// ProcessPriceByExchange - обрабатывает сообщение и возвращает структуру MarketData
func (w *Worker) ProcessFloatsByExchange(msg GenericMessage) storage.MarketData {
	switch data := msg.(type) {
	case BinanceMarketData:
		price, err := strconv.ParseFloat(data.LastPrice, 64)
		if err != nil {
			log.Printf("Failed to parse price: %v", err)
			return storage.MarketData{}
		}

		bid, err := strconv.ParseFloat(data.BestBidPrice, 64)
		if err != nil {
			log.Printf("Failed to parse bid: %v", err)
			return storage.MarketData{}
		}

		ask, err := strconv.ParseFloat(data.BestAskPrice, 64)
		if err != nil {
			log.Printf("Failed to parse ask: %v", err)
			return storage.MarketData{}
		}

		priceInt := int64(price * 1e2)
		bidInt := int64(bid * 1e2)
		askInt := int64(ask * 1e2)

		return storage.MarketData{
			Exchange: "binance",
			Symbol:   data.Symbol,
			Price:    priceInt,
			Bid:      bidInt,
			Ask:      askInt,
		}
	case BybitMarketData:
		price, err := strconv.ParseFloat(data.LastPrice, 64)
		if err != nil {
			log.Printf("Failed to parse price: %v", err)
			return storage.MarketData{}
		}

		priceInt := int64(price * 1e2)

		return storage.MarketData{
			Exchange: "bybit",
			Symbol:   data.Symbol,
			Price:    priceInt,
			Bid:      priceInt,
			Ask:      priceInt,
		}
	default:
		log.Printf("Unsupported type: %T", msg)
		return storage.MarketData{}
	}
}
