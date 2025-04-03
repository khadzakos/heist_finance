package storage

import "time"

type MarketData struct {
	Exchange           string
	Symbol             string
	Market             string
	Price              int64
	Volume             int64
	High               int64
	Low                int64
	PriceChangePercent string
	Timestamp          time.Time
}
