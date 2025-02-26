package app

import (
	"controller/internal/config"
	"controller/internal/controller"
	"log"
	"time"
)

func Run(configPath string) {
	cfg := config.LoadConfig(configPath)

	for _, c := range cfg.Connectors {
		for _, p := range cfg.Preprocessors {
			if c.Queue == p.Queue {
				for err := controller.StartConnectorAndPreprocessor(c, p, cfg.Network); err != nil; {
					log.Println("Ошибка запуска connector и preprocessor", c.Name, p.Name, err)
					time.Sleep(1 * time.Second)
					err = controller.StartConnectorAndPreprocessor(c, p, cfg.Network)
				}
			}
		}
	}

	// go api.StartServer() // TODO: запустить api

	updateChan := make(chan struct{})
	go controller.WatchConfigFile(configPath, updateChan)

	for {
		<-updateChan
		newConfig := config.LoadConfig(configPath)
		controller.UpdateServices(newConfig)
	}
}
