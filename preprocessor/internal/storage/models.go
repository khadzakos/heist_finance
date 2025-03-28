package storage

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	pool *pgxpool.Pool
}

type MarketData struct {
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
	Price    int64  `json:"price"`
	Bid      int64  `json:"bid"`
	Ask      int64  `json:"ask"`
}

type HistoricalData struct {
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
	Open     int64  `json:"open"`
	High     int64  `json:"high"`
	Low      int64  `json:"low"`
	Close    int64  `json:"close"`
	Volume   int64  `json:"volume"`
}
