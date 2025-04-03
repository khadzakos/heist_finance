package service

import "time"

type ResponseMarketData struct {
	Exchange           string
	Symbol             string
	Market             string
	Price              float64
	Volume             float64
	High               float64
	Low                float64
	PriceChangePercent string
	Timestamp          time.Time
}
