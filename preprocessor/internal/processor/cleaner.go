package processor

import (
	"log"
	"preprocessor/internal/storage"
	"time"
)

type Cleaner struct {
	Db *storage.Storage
}

func NewCleaner(db *storage.Storage) (*Cleaner, error) {
	return &Cleaner{
		Db: db,
	}, nil
}

func (c *Cleaner) CleanOldData() {
	log.Println("[CLEANUP] Запуск планировщика очистки")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go c.Db.CleanOldData()

	for range ticker.C {
		go c.Db.CleanOldData()
	}
}
