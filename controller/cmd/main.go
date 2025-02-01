package main

import (
	"controller/internal/api"
	"controller/internal/config"
	"controller/internal/controller"
	"log"
)

func main() {
	// Загружаем конфиг
	cfg := config.LoadConfig("configs/config.yaml")

	// Запускаем все коннекторы из конфига
	for _, conn := range cfg.Connectors {
		go func(conn config.ConnectorConfig) {
			controller.StartConnector(controller.ConnectorConfig{
				Name:    conn.Name,
				Image:   conn.Image,
				WsUrl:   conn.WsUrl,
				Tickers: conn.Tickers,
			})
		}(conn)

	}

	// Запускаем мониторинг
	go controller.MonitorConnectors()

	// Запускаем API для управления коннекторами
	log.Println("Controller started on :8080")
	api.StartServer()

}
