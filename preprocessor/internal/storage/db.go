package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage - структура для работы с БД
type Storage struct {
	pool *pgxpool.Pool
}

// NewStorage - инициализация Storage с пулом подключений
func NewStorage(dbURL string) (*Storage, error) {
	log.Printf("dbURL: %s\n", dbURL)
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &Storage{pool: pool}, nil
}

// Close - закрытие пула соединений
func (s *Storage) Close() {
	s.pool.Close()
}

// ensureTickerExists - добавляет тикер, если он отсутствует
func (s *Storage) ensureTickerExists(ctx context.Context, exchange, symbol string) (int, error) {
	var tickerID int

	// Проверяем, есть ли значение в ENUM
	var exists bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM pg_type t
			JOIN pg_enum e ON t.oid = e.enumtypid
			WHERE t.typname = 'exchange_enum' AND e.enumlabel = $1
		)
	`, exchange).Scan(&exists)

	if err != nil {
		return 0, fmt.Errorf("failed to check enum existence: %w", err)
	}

	// Если биржи нет в ENUM — добавляем её
	if !exists {
		_, err := s.pool.Exec(ctx, fmt.Sprintf(`ALTER TYPE exchange_enum ADD VALUE '%s'`, exchange))
		if err != nil {
			return 0, fmt.Errorf("failed to add new exchange to enum: %w", err)
		}
		log.Printf("✅ Exchange '%s' added to enum", exchange)
	}

	// Теперь можно вставлять тикер
	err = s.pool.QueryRow(ctx, `
		INSERT INTO tickers (exchange, symbol) 
		VALUES ($1, $2)
		ON CONFLICT (exchange, symbol) DO NOTHING
		RETURNING id
	`, exchange, symbol).Scan(&tickerID)

	if err != nil {
		// Если тикер уже добавлен другим процессом, получаем его ID
		err = s.pool.QueryRow(ctx, `
			SELECT id FROM tickers WHERE exchange = $1 AND symbol = $2
		`, exchange, symbol).Scan(&tickerID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert or fetch ticker: %w", err)
		}
	}

	return tickerID, nil
}

// insertMarketData - вставляет данные в market_data
func (s *Storage) insertMarketData(ctx context.Context, tickerID int, price int64) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO market_data (time, ticker_id, price)
		VALUES ($1, $2, $3)
	`, time.Now(), tickerID, price)

	if err != nil {
		return fmt.Errorf("failed to insert market data: %w", err)
	}
	log.Printf("Вставлено в market_data: %d, %d\n", tickerID, price)

	return nil
}

func (s *Storage) SaveTradeData(trade MarketData) error {
	ctx := context.Background()

	log.Printf("⚡ Saving trade data: %+v", trade) // Лог перед вызовом ensureTickerExists

	tickerID, err := s.ensureTickerExists(ctx, trade.Exchange, trade.Ticker)
	if err != nil {
		log.Printf("❌ Failed to ensure ticker exists: %v", err)
		return fmt.Errorf("failed to ensure ticker exists: %w", err)
	}
	log.Printf("✅ TickerID retrieved: %d", tickerID)

	err = s.insertMarketData(ctx, tickerID, trade.Price)
	if err != nil {
		log.Printf("❌ Failed to insert market data: %v", err)
		return fmt.Errorf("failed to insert market data: %w", err)
	}

	log.Printf("✅ Trade data successfully saved!")
	return nil
}

func (s *Storage) TestConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.pool.Ping(ctx)
}
