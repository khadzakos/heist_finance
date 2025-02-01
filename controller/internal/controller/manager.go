package controller

import (
	"controller/internal/config"
	"log"
	"os/exec"
	"sync"
)

var runningConnectors = make(map[string]bool)
var mu sync.Mutex

func StartConnector(c config.Connector) {
	mu.Lock()
	defer mu.Unlock()

	if runningConnectors[c.Name] {
		log.Printf("Коннектор %s уже запущен", c.Name)
		return
	}

	cmd := exec.Command("docker", "run", "-d", "--name", c.Name, c.Image)
	if err := cmd.Run(); err != nil {
		log.Printf("Ошибка запуска %s: %v\n", c.Name, err)
		return
	}

	log.Printf("Запущен коннектор: %s\n", c.Name)
	runningConnectors[c.Name] = true
}

func StopConnector(name string) {
	mu.Lock()
	defer mu.Unlock()

	cmd := exec.Command("docker", "rm", "-f", name)
	if err := cmd.Run(); err != nil {
		log.Printf("Ошибка остановки %s: %v\n", name, err)
		return
	}

	log.Printf("Остановлен коннектор: %s\n", name)
	delete(runningConnectors, name)
}

func RestartConnector(name, image string) {
	StopConnector(name)
	StartConnector(config.Connector{Name: name, Image: image})
}

func UpdateConnectors(newConfig []config.Connector) {
	mu.Lock()
	defer mu.Unlock()

	current := make(map[string]bool)

	// Запускаем новые и обновляем существующие
	for _, c := range newConfig {
		current[c.Name] = true
		if !runningConnectors[c.Name] {
			StartConnector(c)
		}
	}

	// Останавливаем удаленные коннекторы
	for name := range runningConnectors {
		if !current[name] {
			StopConnector(name)
		}
	}
}
