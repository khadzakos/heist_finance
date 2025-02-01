package main

import (
	"controller/internal/api"
	"controller/internal/config"
	"controller/internal/controller"
)

func main() {
	cfg := config.LoadConfig("configs/config.yaml")

	// Запуск коннекторов
	for _, c := range cfg.Connectors {
		controller.StartConnector(c)
	}

	// Запуск API
	go api.StartServer()

	// Мониторинг изменений конфига
	updateChan := make(chan struct{})
	go controller.WatchConfigFile("configs/config.yaml", updateChan)

	for {
		<-updateChan
		newConfig := config.LoadConfig("configs/config.yaml")
		controller.UpdateConnectors(newConfig.Connectors)
	}
}
