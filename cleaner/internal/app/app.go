package app

import (
	"cleaner/internal/cleaner"
	"cleaner/internal/config"
	"cleaner/internal/storage"
	"log"
)

func Run() {
	cfg := config.LoadConfig()

	storage, err := storage.NewStorage(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к DB:", err)
	}
	defer storage.Close()

	cleaner, err := cleaner.NewCleaner(storage)
	if err != nil {
		log.Fatal("Ошибка создания очистителя:", err)
	}

	cleaner.CleanOldData()
}
