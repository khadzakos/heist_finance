package app

import (
	"log"
	"preprocessor/internal/config"
	"preprocessor/internal/processor"
	"preprocessor/internal/storage"
)

func Run() {
	cfg := config.LoadConfig()

	storage, err := storage.NewStorage(cfg.Database.URL)
	if err != nil {
		log.Fatal("Ошибка подключения к DB:", err)
	}
	defer storage.Close()

	err = storage.TestConnection()
	if err != nil {
		log.Fatal("Ошибка тестирования подключения к DB:", err)
	}

	err = processor.ConnectToRabbitMQ(cfg.RabbitMQ.URL)
	defer processor.CloseConnection()
	if err != nil {
		log.Fatal("Ошибка подключения к RabbitMQ:", err)
	}

	processor.ProcessMessages(cfg.Preprocessors, storage)
}
