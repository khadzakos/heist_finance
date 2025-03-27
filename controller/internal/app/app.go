package app

import (
	"controller/internal/config"
	"controller/internal/controller"
	"log"
)

func Run(configPath string) {
	cfg := config.LoadConfig(configPath)

	err := Migrate(config.GetDatabaseURL())
	if err != nil {
		log.Println("Ошибка миграции", err)
	}

	for _, c := range cfg.Connectors {
		for _, p := range cfg.Preprocessors {
			if c.Queue == p.Queue {
				err := controller.StartConnectorAndPreprocessor(c, p, cfg.Network)
				if err != nil {
					log.Println("Ошибка запуска connector и preprocessor", c.Name, p.Name, err)
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
