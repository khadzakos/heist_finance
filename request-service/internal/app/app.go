package app

import (
	"log"
	"net/http"
	"request-service/internal/api"
	"request-service/internal/config"
	"request-service/internal/service"
	"request-service/internal/storage"
)

func Run() {
	cfg := config.LoadConfig()
	db, err := storage.NewStorage(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	marketService := service.NewMarketService(db)

	handler := api.NewHandler(marketService)

	log.Println("Server is running on port 8081")
	http.ListenAndServe(":8081", handler.RegisterRoutes())
}
