package controller

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func WatchConfigFile(path string, updateChan chan struct{}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Ошибка создания watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		log.Fatalf("Ошибка добавления файла в watcher: %v", err)
	}

	log.Println("🔄 Watcher запущен, отслеживание изменений в config.yaml")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Println("📌 Обнаружены изменения в config.yaml, обновляем коннекторы...")
				updateChan <- struct{}{}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Ошибка watcher: %v\n", err)
		}
	}
}
