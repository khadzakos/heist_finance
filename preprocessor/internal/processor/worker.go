package processor

import (
	"log"
	"preprocessor/internal/storage"

	"github.com/streadway/amqp"
)

type Worker struct {
	Id        int
	Jobs      <-chan amqp.Delivery
	Db        *storage.Storage
	Processor *Processor
}

func (w *Worker) ProcessMessages() {
	for msg := range w.Jobs {
		consumedMessage, err := w.Processor.ConsumeMessage(msg.Body)
		if err != nil {
			log.Printf("Worker %d: Ошибка обработки сообщения: %s\n", w.Id, err)
			continue
		}

		processedData := w.ProcessFloatsByExchange(consumedMessage)
		if processedData == (storage.MarketData{}) {
			log.Printf("Worker %d: Не удалось обработать сообщение: %+v\n", w.Id, consumedMessage)
			continue
		}
		w.Db.SaveMarketData(processedData)
		log.Printf("Worker %d: Обработано и сохранено в DB (%s): %+v\n", w.Id, w.Processor.Cfg.Preprocessor.Exchange, processedData)
	}
}
