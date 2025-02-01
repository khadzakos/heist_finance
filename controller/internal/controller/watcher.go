package controller

import (
	"log"
	"os"
	"time"
)

var lastModified time.Time

func WatchConfigFile(path string, updateChan chan struct{}) {
	for {
		info, err := os.Stat(path)
		if err != nil {
			log.Println("Ошибка чтения config.yaml:", err)
			continue
		}

		if info.ModTime().After(lastModified) && !lastModified.IsZero() {
			lastModified = info.ModTime()
			log.Println("Обнаружены изменения в config.yaml, обновляем коннекторы...")
			updateChan <- struct{}{}
		}

		time.Sleep(5 * time.Second)
	}
}
