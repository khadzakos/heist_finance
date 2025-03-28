package processor

import (
	"log"
	"math/rand/v2"
	"preprocessor/internal/config"
	"preprocessor/internal/storage"
	"sync"

	"github.com/streadway/amqp"
)

type Processor struct {
	Cfg            *config.Config
	Db             *storage.Storage
	Conn           *amqp.Connection
	Ch             *amqp.Channel
	workers        []*Worker
	workerChannels []chan amqp.Delivery
	mu             sync.Mutex
	wg             sync.WaitGroup
}

func NewProcessor(cfg *config.Config, db *storage.Storage) (*Processor, error) {
	return &Processor{
		Cfg:            cfg,
		Db:             db,
		Conn:           nil,
		Ch:             nil,
		workers:        make([]*Worker, 0),
		workerChannels: make([]chan amqp.Delivery, 0),
	}, nil
}

func (p *Processor) ProcessMessages(initialWorkerCount int) {
	msgs, err := p.Ch.Consume(
		p.Cfg.Preprocessor.Queue, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Ошибка подписки на очередь:", err)
	}

	for i := 0; i < initialWorkerCount; i++ {
		p.AddWorker()
	}

	const batchSize = 10
	batch := make([]amqp.Delivery, 0, batchSize)

	go func() {
		for msg := range msgs {
			batch = append(batch, msg)
			if len(batch) >= batchSize {
				p.mu.Lock()
				if len(p.workerChannels) > 0 {
					idx := rand.IntN(len(p.workerChannels))
					for _, m := range batch {
						p.workerChannels[idx] <- m
					}
				}
				p.mu.Unlock()
				batch = batch[:0]
			}
		}

		if len(batch) > 0 {
			p.mu.Lock()
			if len(p.workerChannels) > 0 {
				idx := rand.IntN(len(p.workerChannels))
				for _, m := range batch {
					p.workerChannels[idx] <- m
				}
			}
			p.mu.Unlock()
		}
		p.mu.Lock()
		for _, ch := range p.workerChannels {
			close(ch)
		}
		p.mu.Unlock()
	}()

	log.Println("Запуск обработчиков")
	p.wg.Wait()
}

func (p *Processor) AddWorker() *Worker {
	p.mu.Lock()
	defer p.mu.Unlock()

	id := len(p.workers)
	ch := make(chan amqp.Delivery, 100)
	worker := &Worker{
		Id:        id,
		Jobs:      ch,
		Db:        p.Db,
		Processor: p,
		Peers:     make([]*Worker, 0, len(p.workers)),
	}

	for _, w := range p.workers {
		w.Peers = append(w.Peers, worker)
		worker.Peers = append(worker.Peers, w)
	}

	p.workers = append(p.workers, worker)
	p.workerChannels = append(p.workerChannels, ch)

	p.wg.Add(1)
	go func(w *Worker) {
		defer p.wg.Done()
		w.ProcessMessages()
	}(worker)

	log.Printf("Добавлен воркер %d", id)
	return worker
}

func (p *Processor) RemoveWorker(id int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if id < 0 || id >= len(p.workers) {
		return
	}

	close(p.workerChannels[id])

	p.workers = append(p.workers[:id], p.workers[id+1:]...)
	p.workerChannels = append(p.workerChannels[:id], p.workerChannels[id+1:]...)

	for _, w := range p.workers {
		newPeers := make([]*Worker, 0, len(p.workers)-1)
		for _, peer := range w.Peers {
			if peer.Id != id {
				newPeers = append(newPeers, peer)
			}
		}
		w.Peers = newPeers
	}

	log.Printf("Удален воркер %d", id)
}
