package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

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

func NewStorage(dbURL string) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) GetLatestMarketData() ([]MarketData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, `
		SELECT DISTINCT ON (t.exchange, t.symbol, t.market) 
			   t.exchange, 
			   t.symbol, 
			   t.market,
			   m.price, 
			   m.volume,
			   m.high_price,
			   m.low_price,
			   m.price_change_percent,
			   m.timestamp
		FROM market_data m
		JOIN tickers t ON m.ticker_id = t.id
		ORDER BY t.exchange, t.symbol, t.market, m.timestamp DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market data: %w", err)
	}
	defer rows.Close()

	var data []MarketData
	for rows.Next() {
		var d MarketData
		if err := rows.Scan(&d.Exchange, &d.Symbol, &d.Market, &d.Price, &d.Volume, &d.High, &d.Low, &d.PriceChangePercent, &d.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan market data: %w", err)
		}
		data = append(data, d)
	}
	return data, nil
}
