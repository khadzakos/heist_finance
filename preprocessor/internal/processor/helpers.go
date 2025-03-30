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

		volume, err := strconv.ParseFloat(data.TotalTradedBaseAssetVolume, 64)
		if err != nil {
			log.Printf("Failed to parse volume: %v", err)
			return storage.MarketData{}
		}

		high, err := strconv.ParseFloat(data.HighPrice, 64)
		if err != nil {
			log.Printf("Failed to parse high: %v", err)
			return storage.MarketData{}
		}

		low, err := strconv.ParseFloat(data.LowPrice, 64)
		if err != nil {
			log.Printf("Failed to parse low: %v", err)
			return storage.MarketData{}
		}

		priceInt := int64(price * 1e3)
		volumeInt := int64(volume * 1e3)
		highInt := int64(high * 1e3)
		lowInt := int64(low * 1e3)
		if priceInt == 0 {
			return storage.MarketData{}
		}

		return storage.MarketData{
			Exchange:           "binance",
			Symbol:             data.Symbol,
			Price:              priceInt,
			Volume:             volumeInt,
			High:               highInt,
			Low:                lowInt,
			PriceChangePercent: data.PriceChangePercent,
		}
	case BybitMarketData:
		price, err := strconv.ParseFloat(data.LastPrice, 64)
		if err != nil {
			log.Printf("Failed to parse price: %v", err)
			return storage.MarketData{}
		}

		volume, err := strconv.ParseFloat(data.Volume24h, 64)
		if err != nil {
			log.Printf("Failed to parse volume: %v", err)
			return storage.MarketData{}
		}

		high, err := strconv.ParseFloat(data.HighPrice24h, 64)
		if err != nil {
			log.Printf("Failed to parse high: %v", err)
			return storage.MarketData{}
		}

		low, err := strconv.ParseFloat(data.LowPrice24h, 64)
		if err != nil {
			log.Printf("Failed to parse low: %v", err)
			return storage.MarketData{}
		}

		priceInt := int64(price * 1e3)
		volumeInt := int64(volume * 1e3)
		highInt := int64(high * 1e3)
		lowInt := int64(low * 1e3)
		if priceInt == 0 {
			return storage.MarketData{}
		}

		return storage.MarketData{
			Exchange:           "bybit",
			Symbol:             data.Symbol,
			Price:              priceInt,
			Volume:             volumeInt,
			High:               highInt,
			Low:                lowInt,
			PriceChangePercent: data.Price24hPcnt,
		}
	default:
		log.Printf("Unsupported type: %T", msg)
		return storage.MarketData{}
	}
}
