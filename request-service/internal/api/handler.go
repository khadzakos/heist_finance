package api

import (
	"encoding/json"
	"net/http"

	"request-service/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	marketService *service.MarketService
}

func NewHandler(service *service.MarketService) *Handler {
	return &Handler{marketService: service}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", h.GetHomePage).Methods("GET", "OPTIONS")
	r.HandleFunc("/exchange/{exchange}", h.GetExchangeData).Methods("GET", "OPTIONS")
	r.HandleFunc("/exchange/{exchange}/asset/{symbol}", h.GetAssetDetails).Methods("GET", "OPTIONS")

	return corsMiddleware(r)
}

func (h *Handler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	data, err := h.marketService.GetLatestData()
	if err != nil {
		http.Error(w, "Failed to fetch market data", http.StatusInternalServerError)
		return
	}

	crypto := filterByMarket(data, "crypto")
	stock := filterByMarket(data, "stock")

	response := map[string][]service.ResponseMarketData{
		"crypto": crypto,
		"stock":  stock,
	}

	writeJSON(w, response)
}

func (h *Handler) GetCryptoMarket(w http.ResponseWriter, r *http.Request) {
	data, err := h.marketService.GetLatestData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	var filtered []service.ResponseMarketData
	for _, d := range data {
		if d.Market == "crypto" {
			filtered = append(filtered, d)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

func (h *Handler) GetStockMarket(w http.ResponseWriter, r *http.Request) {
	data, err := h.marketService.GetLatestData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	var filtered []service.ResponseMarketData
	for _, d := range data {
		if d.Market == "stock" {
			filtered = append(filtered, d)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

func (h *Handler) GetExchangeData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	exchangeName := vars["exchange"]

	if exchangeName == "" {
		http.Error(w, "Exchange name is required", http.StatusBadRequest)
		return
	}

	data, err := h.marketService.GetLatestData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	var filtered []service.ResponseMarketData
	for _, d := range data {
		if d.Exchange == exchangeName {
			filtered = append(filtered, d)
		}
	}

	if len(filtered) == 0 {
		http.Error(w, "Exchange not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

func (h *Handler) GetAssetDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	exchangeName := vars["exchange"]
	symbol := vars["symbol"]

	if exchangeName == "" || symbol == "" {
		http.Error(w, "Exchange and symbol are required", http.StatusBadRequest)
		return
	}

	data, err := h.marketService.GetLatestData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	for _, d := range data {
		if d.Exchange == exchangeName && d.Symbol == symbol {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(d)
			return
		}
	}

	http.Error(w, "Asset not found", http.StatusNotFound)
}
