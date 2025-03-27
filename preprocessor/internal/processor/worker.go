package processor

import (
	"log"
)

func (w *Worker) ProcessMessages() {
	for msg := range w.Jobs {
		consumedMessage, err := w.Processor.ConsumeMessage(msg.Body)
		if err != nil {
			log.Printf("Worker %d: Ошибка обработки сообщения: %s\n", w.Id, err)
			continue
		}

		processedData := w.ProcessByExchange(w.Processor.Cfg.Preprocessor.Exchange, consumedMessage)
		w.Db.SaveMarketData(processedData)
		log.Printf("Worker %d: Обработано и сохранено в DB (%s): %+v\n", w.Id, w.Processor.Cfg.Preprocessor.Exchange, processedData)
	}
}
