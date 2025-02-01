package api

import (
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/add-connector", AddConnectorHandler)
	http.HandleFunc("/stop-connector", StopConnectorHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
