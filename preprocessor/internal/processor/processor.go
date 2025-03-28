package processor

import (
	"log"
	"preprocessor/internal/config"
	"preprocessor/internal/storage"
	"sync"

	"github.com/streadway/amqp"
)

type Processor struct {
	Cfg  *config.Config
	Db   *storage.Storage
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewProcessor(cfg *config.Config, db *storage.Storage) (*Processor, error) {
	return &Processor{
		Cfg:  cfg,
		Db:   db,
		Conn: nil,
		Ch:   nil,
	}, nil
}

func (p *Processor) ProcessMessages() {
	msgs, err := p.Ch.Consume(
		p.Cfg.Preprocessor.Queue, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Ошибка подписки на очередь:", err)
	}

	var wg sync.WaitGroup
	workers := make([]*Worker, 4)
	for i := 0; i < 4; i++ {
		workers[i] = &Worker{
			Id:        i,
			Jobs:      msgs,
			Db:        p.Db,
			Processor: p,
		}
		wg.Add(1)
		go func(w *Worker) {
			defer wg.Done()
			w.ProcessMessages()
		}(workers[i])
	}

	log.Println("Запуск обработчиков")
	wg.Wait()

}
