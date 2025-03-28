package storage

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ensureTickerExists - добавляет тикер, если он отсутствует, и возвращает его ID
func (s *Storage) ensureTickerExists(ctx context.Context, exchange, symbol string) (int64, error) {
	var tickerID int64

	err := s.pool.QueryRow(ctx, `
		SELECT id FROM tickers 
		WHERE exchange = $1 AND symbol = $2
	`, exchange, symbol).Scan(&tickerID)

	if err == nil {
		return tickerID, nil
	}

	err = s.pool.QueryRow(ctx, `
		INSERT INTO tickers (exchange, symbol) 
		VALUES ($1, $2)
		ON CONFLICT (exchange, symbol) DO NOTHING
		RETURNING id
	`, exchange, symbol).Scan(&tickerID)

	if err != nil {
		err = s.pool.QueryRow(ctx, `
			SELECT id FROM tickers 
			WHERE exchange = $1 AND symbol = $2
		`, exchange, symbol).Scan(&tickerID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert or fetch ticker: %w", err)
		}
	}

	return tickerID, nil
}

// insertMarketData - вставляет текущие рыночные данные
func (s *Storage) insertMarketData(ctx context.Context, tickerID int64, data MarketData) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO market_data (ticker_id, price, bid, ask, timestamp)
		VALUES ($1, $2, $3, $4, $5)
	`, tickerID, data.Price, data.Bid, data.Ask, time.Now().UTC())

	if err != nil {
		return fmt.Errorf("failed to insert market data: %w", err)
	}
	log.Printf("Inserted into market_data: ticker_id=%d, price=%d, bid=%d, ask=%d\n", tickerID, data.Price, data.Bid, data.Ask)

	return nil
}

// insertHistoricalData - вставляет исторические данные
func (s *Storage) insertHistoricalData(ctx context.Context, tickerID int64, data HistoricalData) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO historical_data (ticker_id, open, high, low, close, volume, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (ticker_id, timestamp) DO NOTHING
	`, tickerID, data.Open, data.High, data.Low, data.Close, data.Volume, time.Now().UTC())

	if err != nil {
		return fmt.Errorf("failed to insert historical data: %w", err)
	}
	log.Printf("Inserted into historical_data: ticker_id=%d, open=%d, high=%d, low=%d, close=%d, volume=%d\n",
		tickerID, data.Open, data.High, data.Low, data.Close, data.Volume)

	return nil
}
