package storage

type MarketData struct {
	Exchange  string `json:"exchange"`
	Timestamp int64  `json:"timestamp"`
	Ticker    string `json:"ticker"`
	Price     int64  `json:"price"`
}
