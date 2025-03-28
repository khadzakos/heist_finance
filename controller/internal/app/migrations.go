package app

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate(dbURL string) error {
	var db *sql.DB
	var err error

	maxAttempts := 5
	baseDelay := time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			return fmt.Errorf("error opening database: %w", err)
		}

		err = db.Ping()
		if err == nil {
			break
		}

		delay := baseDelay * time.Duration(math.Pow(2, float64(attempt-1)))

		log.Printf("Attempt %d: Database not ready, waiting %v (Error: %v)",
			attempt, delay, err)

		time.Sleep(delay)

		if attempt == maxAttempts {
			return fmt.Errorf("could not connect to database after %d attempts: %w",
				maxAttempts, err)
		}
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "heist_finance", driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
