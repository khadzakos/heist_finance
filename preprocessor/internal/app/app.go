package app

import (
	"log"
	"preprocessor/internal/config"
	"preprocessor/internal/processor"
	"preprocessor/internal/storage"
)

const initialWorkerCount = 4

func Run() {
	cfg := config.LoadConfig()

	storage, err := storage.NewStorage(cfg.Database.URL)
	if err != nil {
		log.Fatal("Ошибка подключения к DB:", err)
	}
	defer storage.Close()

	p, err := processor.NewProcessor(cfg, storage)
	if err != nil {
		log.Fatal("Ошибка создания процессора:", err)
	}

	c, err := processor.NewCleaner(storage)
	if err != nil {
		log.Fatal("Ошибка создания очистителя:", err)
	}

	err = p.ConnectToRabbitMQ()
	defer p.CloseConnection()
	if err != nil {
		log.Fatal("Ошибка подключения к RabbitMQ:", err)
	}
	log.Println("Подключение к RabbitMQ успешно")

	go p.ProcessMessages(initialWorkerCount)
	go c.CleanOldData()

	select {}
}
