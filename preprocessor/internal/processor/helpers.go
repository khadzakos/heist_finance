package processor

import (
	"log"
	"preprocessor/internal/storage"
	"strconv"
)

// ProcessPriceByExchange - обрабатывает сообщение и возвращает структуру MarketData
func (w *Worker) ProcessFloatsByExchange(msg GenericMessage) storage.MarketData {
	switch data := msg.(type) {
	case BinanceMarketData:
		price, err := strconv.ParseFloat(data.LastPrice, 64)
		if err != nil {
			log.Printf("Failed to parse price: %v", err)
			return storage.MarketData{}
		}

		if price < 0.1 {
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

		return storage.MarketData{
			Exchange:           "binance",
			Symbol:             data.Symbol,
			Market:             "crypto",
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

		if price < 0.1 {
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

		return storage.MarketData{
			Exchange:           "bybit",
			Symbol:             data.Symbol,
			Market:             "crypto",
			Price:              priceInt,
			Volume:             volumeInt,
			High:               highInt,
			Low:                lowInt,
			PriceChangePercent: data.Price24hPcnt,
		}

	case OkxMarketData:
		price, err := strconv.ParseFloat(data.Last, 64)
		if err != nil {
			log.Printf("Failed to parse price: %v", err)
			return storage.MarketData{}
		}

		if price < 0.1 {
			return storage.MarketData{}
		}

		volume, err := strconv.ParseFloat(data.VolCcy24h, 64)
		if err != nil {
			log.Printf("Failed to parse volume: %v", err)
			return storage.MarketData{}
		}

		high, err := strconv.ParseFloat(data.High24h, 64)
		if err != nil {
			log.Printf("Failed to parse high: %v", err)
			return storage.MarketData{}
		}

		low, err := strconv.ParseFloat(data.Low24h, 64)
		if err != nil {
			log.Printf("Failed to parse low: %v", err)
			return storage.MarketData{}
		}

		priceInt := int64(price * 1e3)
		volumeInt := int64(volume * 1e3)
		highInt := int64(high * 1e3)
		lowInt := int64(low * 1e3)

		return storage.MarketData{
			Exchange:           "okx",
			Symbol:             data.InstID,
			Market:             "crypto",
			Price:              priceInt,
			Volume:             volumeInt,
			High:               highInt,
			Low:                lowInt,
			PriceChangePercent: "nil",
		}

	case CoinbaseMarketData:
		price, err := strconv.ParseFloat(data.Price, 64)
		if err != nil {
			log.Printf("Failed to parse price: %v", err)
			return storage.MarketData{}
		}

		if price < 0.1 {
			return storage.MarketData{}
		}

		volume, err := strconv.ParseFloat(data.Volume24h, 64)
		if err != nil {
			log.Printf("Failed to parse volume: %v", err)
			return storage.MarketData{}
		}

		high, err := strconv.ParseFloat(data.High24h, 64)
		if err != nil {
			log.Printf("Failed to parse high: %v", err)
			return storage.MarketData{}
		}

		low, err := strconv.ParseFloat(data.Low24h, 64)
		if err != nil {
			log.Printf("Failed to parse low: %v", err)
			return storage.MarketData{}
		}

		priceInt := int64(price * 1e3)
		volumeInt := int64(volume * 1e3)
		highInt := int64(high * 1e3)
		lowInt := int64(low * 1e3)

		return storage.MarketData{
			Exchange:           "coinbase",
			Symbol:             data.ProductID,
			Market:             "crypto",
			Price:              priceInt,
			Volume:             volumeInt,
			High:               highInt,
			Low:                lowInt,
			PriceChangePercent: "nil",
		}

	default:
		log.Printf("Unsupported type: %T", msg)
		return storage.MarketData{}
	}
}
