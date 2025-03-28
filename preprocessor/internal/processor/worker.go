package processor

import (
	"log"
	"preprocessor/internal/storage"
	"time"

	"github.com/streadway/amqp"
)

type Worker struct {
	Id        int
	Jobs      chan amqp.Delivery
	Db        *storage.Storage
	Processor *Processor
	Peers     []*Worker
}

func (w *Worker) ProcessMessages() {
	for {
		select {
		case msg, ok := <-w.Jobs:
			if !ok {
				log.Printf("Worker %d: Канал закрыт, завершение работы", w.Id)
				return
			}
			w.processMessage(msg)
		default:
			stole := false
			for _, peer := range w.Peers {
				select {
				case msg, ok := <-peer.Jobs:
					if ok {
						log.Printf("Worker %d: Украл задачу у Worker %d", w.Id, peer.Id)
						w.processMessage(msg)
						stole = true
					}
				default:
					continue
				}

			}
			if !stole {
				time.Sleep(10 * time.Second)
			}
		}
	}
}

func (w *Worker) processMessage(msg amqp.Delivery) {
	consumedMessage, err := w.Processor.ConsumeMessage(msg.Body)
	if err != nil {
		log.Printf("Worker %d: Ошибка обработки сообщения: %s", w.Id, err)
		return
	}

	processedData := w.ProcessFloatsByExchange(consumedMessage)
	if processedData == (storage.MarketData{}) {
		log.Printf("Worker %d: Не удалось обработать сообщение: %+v", w.Id, consumedMessage)
		return
	}
	err = w.Db.SaveMarketData(processedData)
	if err != nil {
		log.Printf("Worker %d: Ошибка сохранения данных: %s", w.Id, err)
		return
	}
	log.Printf("Worker %d: Обработано и сохранено в DB (%s): %+v", w.Id, w.Processor.Cfg.Preprocessor.Exchange, processedData)
}
