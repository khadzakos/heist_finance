package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	lockID = int64(12345)
)

type Storage struct {
	pool *pgxpool.Pool
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

func (s *Storage) CleanOldData() {
	ctx := context.Background()
	log.Println("Попытка получить advisory lock...")

	var acquired bool
	err := s.pool.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", lockID).Scan(&acquired)
	if err != nil || !acquired {
		log.Println("Lock занят другим инстансом, пропуск.")
		return
	}
	defer func() {
		_, _ = s.pool.Exec(ctx, "SELECT pg_advisory_unlock($1)", lockID)
	}()

	var count int
	err = s.pool.QueryRow(ctx, `
        SELECT COUNT(*) FROM market_data 
        WHERE timestamp < NOW() - INTERVAL '5 minute'
    `).Scan(&count)
	if err != nil {
		log.Printf("Ошибка при подсчёте строк: %v", err)
		return
	}
	log.Printf("К удалению найдено %d строк (старше 5 минут)", count)

	res, err := s.pool.Exec(ctx, `
        DELETE FROM market_data 
        WHERE timestamp < NOW() - INTERVAL '5 minute'
    `)
	if err != nil {
		log.Printf("Ошибка при удалении: %v", err)
		return
	}

	log.Printf("Удалено %d строк", res.RowsAffected())
}
