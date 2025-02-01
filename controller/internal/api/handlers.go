package api

import (
	"controller/internal/config"
	"controller/internal/controller"
	"encoding/json"
	"net/http"
)

func AddConnectorHandler(w http.ResponseWriter, r *http.Request) {
	var cfg config.Connector
	json.NewDecoder(r.Body).Decode(&cfg)

	go controller.StartConnector(cfg)
	w.Write([]byte("Коннектор запущен"))
}

func StopConnectorHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&request)

	controller.StopConnector(request.Name)
	w.Write([]byte("Коннектор остановлен"))
}
