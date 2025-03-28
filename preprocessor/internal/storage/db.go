package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (s *Storage) SaveMarketData(data MarketData) error {
	ctx := context.Background()

	tickerID, err := s.ensureTickerExists(ctx, data.Exchange, data.Symbol)
	if err != nil {
		return fmt.Errorf("failed to ensure ticker exists: %w", err)
	}

	err = s.insertMarketData(ctx, tickerID, data)
	if err != nil {
		return fmt.Errorf("failed to insert market data: %w", err)
	}
	log.Printf("Market data inserted: %+v", data)
	return nil
}

func (s *Storage) SaveHistoricalData(data HistoricalData) error {
	ctx := context.Background()

	tickerID, err := s.ensureTickerExists(ctx, data.Exchange, data.Symbol)
	if err != nil {
		return fmt.Errorf("failed to ensure ticker exists: %w", err)
	}

	err = s.insertHistoricalData(ctx, tickerID, data)
	if err != nil {
		return fmt.Errorf("failed to insert historical data: %w", err)
	}

	return nil
}
