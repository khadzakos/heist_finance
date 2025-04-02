package api

import (
	"encoding/json"
	"net/http"

	"request-service/internal/service"
	"request-service/internal/storage"

	"github.com/gorilla/mux"
)

type Handler struct {
	marketService *service.MarketService
}

func NewHandler(service *service.MarketService) *Handler {
	return &Handler{marketService: service}
}

func (h *Handler) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", h.GetHomePage).Methods("GET")
	r.HandleFunc("/crypto-market", h.GetCryptoMarket).Methods("GET")
	// r.HandleFunc("/stock-market", h.GetStockMarket).Methods("GET")
	// r.HandleFunc("/exchange/{exchange}", h.GetExchangeData).Methods("GET")
	// r.HandleFunc("/asset/{exchange}/{symbol}", h.GetAssetDetails).Methods("GET")

	return r
}

func (h *Handler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	data, err := h.marketService.GetLatestData()
	if err != nil {
		http.Error(w, "Failed to fetch market data", http.StatusInternalServerError)
		return
	}

	crypto := filterByMarket(data, "crypto")
	stock := filterByMarket(data, "stock")

	response := map[string][]storage.MarketData{
		"crypto": topAssets(crypto, 5),
		"stock":  topAssets(stock, 5),
	}

	writeJSON(w, response)
}

func (h *Handler) GetCryptoMarket(w http.ResponseWriter, r *http.Request) {
	data, err := h.marketService.GetLatestData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	var filtered []storage.MarketData
	for _, d := range data {
		if d.Market == "crypto" {
			filtered = append(filtered, d)
		}
	}

	json.NewEncoder(w).Encode(filtered)
}
