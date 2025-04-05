package cleaner

import (
	"cleaner/internal/storage"
	"log"
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
	log.Println("Запуск планировщика очистки")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	c.Db.CleanOldData()

	for range ticker.C {
		c.Db.CleanOldData()
	}
}
